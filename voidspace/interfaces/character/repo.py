from typing import Protocol
from uuid import UUID

from voidspace.database.models import Character


class CharacterRepo(Protocol):
    def create_character(self, character: Character) -> None:
        raise NotImplementedError

    async def find_one_character(self, id: UUID) -> Character | None:
        raise NotImplementedError

    async def find_all_characters_by_user_id(self, user_id: UUID) -> list[Character]:
        raise NotImplementedError

    async def find_one_character_by_code(self, code: str) -> Character | None:
        raise NotImplementedError

    async def delete_character(self, id: UUID) -> Character | None:
        raise NotImplementedError

    def update_character(self, character: Character) -> None:
        raise NotImplementedError
