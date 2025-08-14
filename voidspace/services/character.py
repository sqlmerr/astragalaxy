from dataclasses import dataclass
from uuid import UUID

from voidspace.database.models import Character
from voidspace.dto.character import CharacterDTO, CreateCharacterDTO
from voidspace.exceptions.character import (
    CharacterNotFound,
    CharacterCodeAlreadyOccupied,
)
from voidspace.interfaces.character.create import CharacterWriter
from voidspace.interfaces.character.delete import CharacterDeleter
from voidspace.interfaces.character.read import CharacterReader
from voidspace.interfaces.character.repo import CharacterRepo


@dataclass(frozen=True)
class CharacterService(CharacterReader, CharacterWriter, CharacterDeleter):
    repo: CharacterRepo

    async def get_character_by_id(self, id: UUID) -> CharacterDTO:
        character = await self.repo.find_one_character(id)
        if not character:
            raise CharacterNotFound

        return CharacterDTO.from_model(character)

    async def get_character_by_code(self, code: str) -> CharacterDTO:
        character = await self.repo.find_one_character_by_code(code)
        if not character:
            raise CharacterNotFound

        return CharacterDTO.from_model(character)

    async def get_characters_by_user(self, user_id: UUID) -> list[CharacterDTO]:
        characters = await self.repo.find_all_characters_by_user_id(user_id)
        return [CharacterDTO.from_model(character) for character in characters]

    async def create_character(self, data: CreateCharacterDTO) -> CharacterDTO:
        character = await self.repo.find_one_character_by_code(data.code.lower())
        if character:
            raise CharacterCodeAlreadyOccupied

        c = Character(
            code=data.code.lower(), location="space_station", user_id=data.user_id
        )
        character_id = await self.repo.create_character(c)

        return await self.get_character_by_id(character_id)

    async def delete_character(self, id: UUID) -> None:
        character = await self.repo.find_one_character(id)
        if not character:
            raise CharacterNotFound

        await self.repo.delete_character(id)
