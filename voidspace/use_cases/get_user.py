from dataclasses import dataclass
from uuid import UUID

from voidspace.dto.user import UserDTO
from voidspace.exceptions import AccessDeniedError
from voidspace.exceptions.user import UserNotFound
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.user import UserRepo
from voidspace.use_cases import BaseUseCase


@dataclass(frozen=True)
class GetUserById(BaseUseCase[UUID, UserDTO]):
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
class GetUserByUsername(BaseUseCase[str, UserDTO]):
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
