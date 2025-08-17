from dataclasses import dataclass
from uuid import uuid4

from voidspace.database.models import Character, Spaceship
from voidspace.dto.character import CreateCharacterDTO, CharacterDTO
from voidspace.exceptions import AppError
from voidspace.exceptions.character import (
    TooManyCharacters,
    CharacterCodeAlreadyOccupied,
)
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.character.repo import CharacterRepo
from voidspace.interfaces.spaceship.repo import SpaceshipRepo
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.use_cases import BaseUseCase


@dataclass(frozen=True)
class CreateCharacter(BaseUseCase[CreateCharacterDTO, CharacterDTO]):
    character_repo: CharacterRepo
    spaceship_repo: SpaceshipRepo
    system_repo: SystemRepo
    idp: IdentityProvider

    async def execute(self, data: CreateCharacterDTO) -> CharacterDTO:
        current_user_id = self.idp.get_current_user_id()
        user_characters = await self.character_repo.find_all_characters_by_user_id(
            current_user_id
        )
        if len(user_characters) >= 5:
            raise TooManyCharacters()

        character = await self.character_repo.find_one_character_by_code(
            data.code.lower()
        )
        if character:
            raise CharacterCodeAlreadyOccupied()

        system = await self.system_repo.find_one_random_system()
        if not system:
            raise AppError

        c = Character(
            id=uuid4(),
            code=data.code.lower(),
            location="space_station",
            user_id=current_user_id,
            system_id=system.id,
        )
        self.character_repo.create_character(c)

        sp = Spaceship(
            id=uuid4(),
            name="initial",
            location="space_station",
            character_id=c.id,
            active=True,
            system_id=system.id,
            planet_id=None,
        )
        self.spaceship_repo.add_spaceship(sp)

        character = await self.character_repo.find_one_character(c.id)

        return CharacterDTO.from_model(character)
