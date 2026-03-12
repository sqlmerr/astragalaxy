from dataclasses import dataclass

from astragalaxy.dto.user import CreateUserDTO, UserDTO
from astragalaxy.exceptions import AppError
from astragalaxy.interfaces.user import UserRepo
from astragalaxy.password_hasher import PasswordHasher
from ..exceptions.user import UsernameAlreadyOccupied
from ..interfaces.user.repo import UserCreate
from ..utils import generate_user_token


@dataclass(frozen=True)
class Register:
    user_repo: UserRepo
    password_hasher: PasswordHasher

    async def execute(self, data: CreateUserDTO) -> UserDTO:
        user = await self.user_repo.find_one_user_by_username(data.username.lower())
        if user:
            raise UsernameAlreadyOccupied()

        hashed_password = self.password_hasher.hash_password(data.password)
        user_id = await self.user_repo.create_user(
            UserCreate(
                data.username.lower(),
                password=hashed_password,
                token=generate_user_token(),
            )
        )
        user = await self.user_repo.find_one_user(user_id)
        if not user:
            raise AppError()

        return UserDTO.from_model(user)
