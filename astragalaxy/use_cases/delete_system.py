from dataclasses import dataclass

from astragalaxy.exceptions import AccessDeniedError
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.system.repo import SystemRepo


@dataclass(frozen=True)
class DeleteSystem:
    repo: SystemRepo
    idp: IdentityProvider
    commiter: Commiter

    async def execute(self, data: str) -> None:
        current_user = await self.idp.get_current_user()
        if current_user.username != "admin": # TODO: roles
            raise AccessDeniedError()

        await self.repo.delete_system(data)
        await self.commiter.commit()