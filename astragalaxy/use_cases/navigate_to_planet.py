import random
from dataclasses import dataclass

from astragalaxy.cooldown_manager import CooldownManager
from astragalaxy.dto.cooldown import CooldownDTO, SetCooldownDTO
from astragalaxy.exceptions import CharacterInCooldown, AppError
from astragalaxy.exceptions.planet import PlanetNotFound
from astragalaxy.exceptions.spaceship import CharacterNeedsToBeInSpaceship
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo
from astragalaxy.interfaces.planet.repo import PlanetRepo
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class NavigateToPlanet:
    spaceship_repo: SpaceshipRepo
    character_repo: CharacterRepo
    planet_repo: PlanetRepo
    idp: IdentityProvider
    cooldown_manager: CooldownManager

    async def execute(self, data: str) -> CooldownDTO:
        current_character_id = self.idp.get_current_character_id()
        current_character = await self.character_repo.find_one_character(
            current_character_id
        )
        if not current_character:
            raise AppError()

        spaceship = await self.spaceship_repo.find_one_active_by_character(
            current_character_id
        )
        if not spaceship:
            raise AppError()

        cooldown = await self.cooldown_manager.get(current_character_id)
        if cooldown.remaining_seconds > 0:
            raise CharacterInCooldown()

        planet = await self.planet_repo.find_one_planet(data)
        if planet.system_id != current_character.system_id:
            raise PlanetNotFound()

        if not current_character.in_spaceship:
            raise CharacterNeedsToBeInSpaceship()

        seconds = random.randint(2 * 60, 3 * 60)
        cooldown = await self.cooldown_manager.set(
            SetCooldownDTO(
                character_id=current_character_id,
                seconds=seconds,
                action="navigation_to_planet",
            )
        )
        current_character.planet_id = data
        current_character.location = "planet"
        spaceship.planet_id = data
        spaceship.location = "planet"
        self.spaceship_repo.save_spaceship(spaceship)
        self.character_repo.update_character(current_character)

        return cooldown
