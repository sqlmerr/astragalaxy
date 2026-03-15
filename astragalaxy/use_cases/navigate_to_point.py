import random
from dataclasses import dataclass

from astragalaxy.cooldown_manager import CooldownManager
from astragalaxy.dto.cooldown import CooldownDTO, SetCooldownDTO
from astragalaxy.exceptions import CharacterInCooldown, AppError
from astragalaxy.exceptions.character import CharacterNotFound
from astragalaxy.exceptions.planet import PlanetNotFound
from astragalaxy.exceptions.point import PointNotFound
from astragalaxy.exceptions.spaceship import CharacterNeedsToBeInSpaceship
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo
from astragalaxy.interfaces.planet.repo import PlanetRepo
from astragalaxy.interfaces.point.repo import PointRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class NavigateToPoint:
    spaceship_repo: SpaceshipRepo
    character_repo: CharacterRepo
    planet_repo: PlanetRepo
    point_repo: PointRepo
    idp: IdentityProvider
    cooldown_manager: CooldownManager
    commiter: Commiter

    async def execute(self, point_id: str) -> CooldownDTO:
        current_character_id = self.idp.get_current_character_id()
        current_character = await self.character_repo.find_one_character(
            current_character_id
        )
        if not current_character:
            raise CharacterNotFound()

        spaceship = await self.spaceship_repo.find_one_active_by_character(
            current_character_id
        )
        if not spaceship:
            raise AppError()

        cooldown = await self.cooldown_manager.get(current_character_id)
        if cooldown.remaining_seconds > 0:
            raise CharacterInCooldown()

        current_point = await self.point_repo.find_one_point(current_character.point_id)
        if not current_point:
            raise AppError()

        point = await self.point_repo.find_one_point(point_id)
        if not point:
            raise PointNotFound()
        

        if point.system_id != current_point.system_id:
            raise PointNotFound()

        if not current_character.in_spaceship:
            raise CharacterNeedsToBeInSpaceship()

        seconds = random.randint(2 * 60, 3 * 60)
        cooldown = await self.cooldown_manager.set(
            SetCooldownDTO(
                character_id=current_character_id,
                seconds=seconds,
                action="navigation_to_point",
            )
        )
        current_character.point_id = point_id
        spaceship.point_id = point_id
        self.spaceship_repo.add(spaceship)
        self.character_repo.add(current_character)
        await self.commiter.commit()

        return cooldown
