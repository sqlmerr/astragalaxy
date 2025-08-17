from dataclasses import dataclass
from uuid import UUID

from voidspace.cooldown_manager import CooldownManager
from voidspace.dto.cooldown import SetCooldownDTO
from voidspace.exceptions import CharacterInCooldown, AppError, AccessDeniedError
from voidspace.exceptions.spaceship import (
    CharacterAlreadyInSpaceship,
    SpaceshipNotFoundError,
)
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.character.repo import CharacterRepo
from voidspace.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class EnterSpaceship:
    repo: SpaceshipRepo
    character_repo: CharacterRepo
    cooldown_manager: CooldownManager
    idp: IdentityProvider

    async def execute(self, data: UUID) -> None:
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

        if active_spaceship.id != data:
            spaceship = await self.repo.find_one_by_id(data)
            if not spaceship:
                raise SpaceshipNotFoundError()
            if spaceship.character_id != current_character_id:
                raise AccessDeniedError()

            spaceship.active = True
            active_spaceship.active = False
            self.repo.add_spaceship(spaceship)
            self.repo.add_spaceship(active_spaceship)

        current_character.in_spaceship = True
        self.character_repo.update_character(current_character)
        await self.cooldown_manager.set(
            SetCooldownDTO(
                character_id=current_character_id, seconds=3, action="spaceship_enter"
            )
        )
