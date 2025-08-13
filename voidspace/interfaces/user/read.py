from uuid import UUID

from typing import Protocol

from voidspace.dto.user import UserDTO


class UserReader(Protocol):
    async def get_user_by_id(self, id: UUID) -> UserDTO:
        raise NotImplementedError

    async def get_user_by_username(self, username: str) -> UserDTO:
        raise NotImplementedError
