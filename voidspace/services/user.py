from dataclasses import dataclass
from uuid import UUID

from voidspace.dto.user import CreateUserDTO, UserDTO, LoginUserDTO
from voidspace.exceptions.user import (
    UsernameAlreadyOccupied,
    UserNotFound,
    InvalidCredentials,
)
from voidspace.interfaces.user import UserReader, UserWriter, UserRepo
from voidspace.interfaces.user.repo import UserCreate
from voidspace.jwt_token_processor import JwtTokenProcessor
from voidspace.password_hasher import PasswordHasher
from voidspace.utils import generate_user_token


@dataclass(frozen=True)
class UserService(UserReader, UserWriter):
    repo: UserRepo
    password_hasher: PasswordHasher
    jwt_token_processor: JwtTokenProcessor

    async def get_user_by_id(self, id: UUID) -> UserDTO:
        user = await self.repo.find_one_user(id)
        if not user:
            raise UserNotFound

        return UserDTO.from_model(user)

    async def get_user_by_username(self, username: str) -> UserDTO:
        user = await self.repo.find_one_user_by_username(username.lower())
        if not user:
            raise UserNotFound

        return UserDTO.from_model(user)

    async def login(self, data: LoginUserDTO) -> str:
        user = await self.repo.find_one_user_by_username(data.username.lower())
        if not user:
            raise InvalidCredentials

        if not self.password_hasher.verify_password(data.password, user.password):
            raise InvalidCredentials

        token = self.jwt_token_processor.encode({"sub": str(user.id)})

        return token

    async def create_user(self, data: CreateUserDTO) -> UserDTO:
        user = await self.repo.find_one_user_by_username(data.username.lower())
        if user:
            raise UsernameAlreadyOccupied

        hashed_password = self.password_hasher.hash_password(data.password)
        user_id = await self.repo.create_user(
            UserCreate(
                data.username.lower(),
                password=hashed_password,
                token=generate_user_token(),
            )
        )
        return await self.get_user_by_id(user_id)
