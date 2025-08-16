from dataclasses import dataclass

from voidspace.exceptions import AccessDeniedError
from voidspace.exceptions.planet import PlanetNotFound
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.planet.repo import PlanetRepo
from voidspace.use_cases import BaseUseCase


@dataclass(frozen=True)
class DeletePlanet(BaseUseCase[str, None]):
    repo: PlanetRepo
    idp: IdentityProvider

    async def execute(self, data: str) -> None:
        current_user = await self.idp.get_current_user()
        if current_user.username != "admin":
            raise AccessDeniedError()

        planet = await self.repo.find_one_planet(data)
        if not planet:
            raise PlanetNotFound()

        await self.repo.delete_planet(data)
