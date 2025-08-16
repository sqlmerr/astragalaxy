from dataclasses import dataclass

from voidspace.exceptions import AccessDeniedError
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.use_cases import BaseUseCase


@dataclass(frozen=True)
class DeleteSystem(BaseUseCase[str, None]):
    repo: SystemRepo
    idp: IdentityProvider

    async def execute(self, data: str) -> None:
        current_user = await self.idp.get_current_user()
        if current_user.username != "admin":
            raise AccessDeniedError()

        await self.repo.delete_system(data)
