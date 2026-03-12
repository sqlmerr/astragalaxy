from dataclasses import dataclass

from astragalaxy.dto.user import LoginUserDTO, AuthTokenDTO
from astragalaxy.exceptions.user import InvalidCredentials
from astragalaxy.interfaces.user.repo import UserRepo
from astragalaxy.jwt_token_processor import JwtTokenProcessor
from astragalaxy.password_hasher import PasswordHasher


@dataclass(frozen=True)
class Login:
    user_repo: UserRepo
    password_hasher: PasswordHasher
    jwt_token_processor: JwtTokenProcessor

    async def execute(self, data: LoginUserDTO) -> AuthTokenDTO:
        user = await self.user_repo.find_one_user_by_username(data.username.lower())
        if not user:
            raise InvalidCredentials()

        if not self.password_hasher.verify_password(data.password, user.password):
            raise InvalidCredentials()

        token = self.jwt_token_processor.encode({"sub": str(user.id)})

        return AuthTokenDTO(access_token=token)
