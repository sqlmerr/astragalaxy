from typing import Protocol
from uuid import UUID


class CharacterDeleter(Protocol):
    async def delete_character(self, id: UUID) -> None:
        raise NotImplementedError
