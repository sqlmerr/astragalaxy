from dataclasses import dataclass

from astragalaxy.exceptions import AccessDeniedError
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.system.repo import SystemRepo


@dataclass(frozen=True)
class DeleteSystem:
    repo: SystemRepo
    idp: IdentityProvider

    async def execute(self, data: str) -> None:
        current_user = await self.idp.get_current_user()
        if current_user.username != "admin":
            raise AccessDeniedError()

        await self.repo.delete_system(data)
