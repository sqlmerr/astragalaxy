from dataclasses import dataclass

from astragalaxy.dto.spaceship import RenameSpaceshipDTO, SpaceshipDTO
from astragalaxy.exceptions import AccessDeniedError, AppError
from astragalaxy.exceptions.spaceship import (
    SpaceshipNotFoundError,
    InvalidSpaceshipNameError,
)
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class RenameSpaceship:
    repo: SpaceshipRepo
    character_repo: CharacterRepo
    idp: IdentityProvider
    commiter: Commiter

    async def execute(self, data: RenameSpaceshipDTO) -> SpaceshipDTO:
        character = await self.idp.get_current_character()
        current_user = await self.idp.get_current_user()

        spaceship = await self.repo.find_one_by_id(data.spaceship_id)
        if not spaceship:
            raise SpaceshipNotFoundError()

        if spaceship.character_id != character.id:
            spaceship_character = await self.character_repo.find_one_character(
                spaceship.character_id
            )
            if not spaceship_character:
                raise SpaceshipNotFoundError()
            if spaceship_character.user_id != current_user.id:
                raise AccessDeniedError()

        if len(data.name) >= 32 or len(data.name) < 3:
            raise InvalidSpaceshipNameError()

        spaceship.name = data.name
        self.repo.add(spaceship)

        await self.commiter.commit()

        dto = SpaceshipDTO.from_model(spaceship)
        return dto
