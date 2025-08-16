from dataclasses import dataclass

from voidspace.dto.spaceship import RenameSpaceshipDTO, SpaceshipDTO
from voidspace.exceptions import AppError
from voidspace.exceptions.spaceship import (
    SpaceshipNotFoundError,
    InvalidSpaceshipNameError,
)
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.character.repo import CharacterRepo
from voidspace.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class RenameSpaceship:
    repo: SpaceshipRepo
    character_repo: CharacterRepo
    idp: IdentityProvider

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
            if not character:
                raise AppError()
            if spaceship_character.user_id != current_user.id:
                raise SpaceshipNotFoundError()

        if len(data.name) >= 32 or len(data.name) < 3:
            raise InvalidSpaceshipNameError()

        spaceship.name = data.name
        self.repo.save_spaceship(spaceship)

        return SpaceshipDTO.from_model(spaceship)
