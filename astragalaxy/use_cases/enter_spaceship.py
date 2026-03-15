from dataclasses import dataclass
from uuid import UUID

from astragalaxy.cooldown_manager import CooldownManager
from astragalaxy.dto.cooldown import SetCooldownDTO
from astragalaxy.exceptions import CharacterInCooldown, AppError, AccessDeniedError
from astragalaxy.exceptions.spaceship import (
    CharacterAlreadyInSpaceship,
    SpaceshipAtDifferentPoint,
    SpaceshipNotFoundError,
)
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class EnterSpaceship:
    repo: SpaceshipRepo
    character_repo: CharacterRepo
    cooldown_manager: CooldownManager
    idp: IdentityProvider
    commiter: Commiter

    async def execute(self, spaceship_id: UUID) -> None:
        current_character_id = self.idp.get_current_character_id()
        current_character = await self.character_repo.find_one_character(
            current_character_id
        )
        if not current_character:
            raise AppError()
        cooldown = await self.cooldown_manager.get(current_character.id)
        if cooldown.remaining_seconds > 0:
            raise CharacterInCooldown()

        active_spaceship = await self.repo.find_one_active_by_character(
            current_character.id
        )
        if not active_spaceship:
            raise AppError()

        if current_character.in_spaceship:
            raise CharacterAlreadyInSpaceship()

        if active_spaceship.id != spaceship_id:
            spaceship = await self.repo.find_one_by_id(spaceship_id)
            if not spaceship:
                raise SpaceshipNotFoundError()
            if spaceship.character_id != current_character_id:
                raise AccessDeniedError()
            
            if spaceship.point_id != current_character.point_id:
                raise SpaceshipAtDifferentPoint()

            spaceship.active = True
            active_spaceship.active = False
            self.repo.add(spaceship)
            self.repo.add(active_spaceship)

        current_character.in_spaceship = True
        self.character_repo.add(current_character)
        await self.cooldown_manager.set(
            SetCooldownDTO(
                character_id=current_character_id, seconds=3, action="spaceship_enter"
            )
        )
        await self.commiter.commit()
