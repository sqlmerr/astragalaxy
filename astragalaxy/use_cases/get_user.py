from dataclasses import dataclass
from uuid import UUID

from astragalaxy.dto.user import UserDTO
from astragalaxy.exceptions import AccessDeniedError
from astragalaxy.exceptions.user import UserNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.user import UserRepo


@dataclass(frozen=True)
class GetUserById:
    user_repo: UserRepo
    idp: IdentityProvider

    async def execute(self, data: UUID) -> UserDTO:
        user_id = self.idp.get_current_user_id()
        if data != user_id:
            raise AccessDeniedError()

        user = await self.user_repo.find_one_user(data)
        if not user:
            raise UserNotFound()

        return UserDTO.from_model(user)


@dataclass(frozen=True)
class GetUserByUsername:
    user_repo: UserRepo
    idp: IdentityProvider

    async def execute(self, data: str) -> UserDTO:
        user_id = self.idp.get_current_user_id()
        if data != user_id:
            raise AccessDeniedError()

        user = await self.user_repo.find_one_user_by_username(data)
        if not user:
            raise UserNotFound()

        return UserDTO.from_model(user)
