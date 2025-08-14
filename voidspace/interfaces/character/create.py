from typing import Protocol

from voidspace.dto.character import CharacterDTO, CreateCharacterDTO


class CharacterWriter(Protocol):
    async def create_character(self, data: CreateCharacterDTO) -> CharacterDTO:
        raise NotImplementedError
