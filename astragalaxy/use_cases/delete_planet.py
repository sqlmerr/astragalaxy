from dataclasses import dataclass

from astragalaxy.exceptions import AccessDeniedError
from astragalaxy.exceptions.planet import PlanetNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.planet.repo import PlanetRepo


@dataclass(frozen=True)
class DeletePlanet:
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
