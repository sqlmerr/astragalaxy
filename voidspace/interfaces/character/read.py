from typing import Protocol
from uuid import UUID

from voidspace.dto.character import CharacterDTO


class CharacterReader(Protocol):
    async def get_character_by_id(self, id: UUID) -> CharacterDTO:
        raise NotImplementedError

    async def get_character_by_code(self, code: str) -> CharacterDTO:
        raise NotImplementedError

    async def get_characters_by_user(self, user_id: UUID) -> list[CharacterDTO]:
        raise NotImplementedError
