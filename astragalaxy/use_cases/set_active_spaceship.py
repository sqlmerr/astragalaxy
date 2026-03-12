from dataclasses import dataclass
from uuid import UUID

from astragalaxy.dto.spaceship import SpaceshipDTO
from astragalaxy.exceptions.spaceship import SpaceshipNotFoundError
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class SetActiveSpaceship:
    repo: SpaceshipRepo
    idp: IdentityProvider

    async def execute(self, data: UUID) -> SpaceshipDTO:
        current_character = await self.idp.get_current_character()
        spaceships = await self.repo.find_all_by_character_id(current_character.id)

        if len(spaceships) == 0:
            raise SpaceshipNotFoundError()

        active_spaceship = await self.repo.find_one_active_by_character(
            current_character.id
        )

        for sp in spaceships:
            if sp.id != data:
                continue
            if active_spaceship:
                active_spaceship.active = False
                self.repo.save_spaceship(active_spaceship)
            sp.active = True
            self.repo.save_spaceship(sp)
            return SpaceshipDTO.from_model(sp)

        raise SpaceshipNotFoundError()
