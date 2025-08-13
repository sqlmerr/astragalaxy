from typing import Protocol

from voidspace.dto.user import CreateUserDTO, UserDTO


class UserWriter(Protocol):
    async def create_user(self, data: CreateUserDTO) -> UserDTO:
        raise NotImplementedError
