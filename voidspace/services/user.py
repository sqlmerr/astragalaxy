from uuid import UUID
from dataclasses import dataclass

from voidspace.exceptions.user import UsernameAlreadyOccupied, UserNotFound
from voidspace.interfaces.user import UserReader, UserWriter, UserRepo
from voidspace.dto.user import CreateUserDTO, UserDTO
from voidspace.interfaces.user.repo import UserCreate
from voidspace.utils import generate_user_token


@dataclass(frozen=True)
class UserService(UserReader, UserWriter):
    repo: UserRepo

    async def get_user_by_id(self, id: UUID) -> UserDTO:
        user = await self.repo.find_one_user(id)
        if not user:
            raise UserNotFound

        return UserDTO.from_model(user)

    async def get_user_by_username(self, username: str) -> UserDTO:
        user = await self.repo.find_one_user_by_username(username)
        if not user:
            raise UserNotFound

        return UserDTO.from_model(user)

    async def create_user(self, data: CreateUserDTO) -> UserDTO:
        user = await self.repo.find_one_user_by_username(data.username.lower())
        if user:
            raise UsernameAlreadyOccupied

        user_id = await self.repo.create_user(
            UserCreate(
                data.username, password=data.password, token=generate_user_token()
            )
        )
        return await self.get_user_by_id(user_id)
