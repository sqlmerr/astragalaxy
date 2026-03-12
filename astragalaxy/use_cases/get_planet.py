from dataclasses import dataclass

from astragalaxy.dto.planet import PlanetDTO
from astragalaxy.exceptions.planet import PlanetNotFound
from astragalaxy.exceptions.system import SystemNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.planet.repo import PlanetRepo
from astragalaxy.interfaces.system.repo import SystemRepo


@dataclass(frozen=True)
class GetPlanet:
    repo: PlanetRepo
    idp: IdentityProvider

    async def execute(self, data: str) -> PlanetDTO:
        await self.idp.get_current_user()

        planet = await self.repo.find_one_planet(data)
        if not planet:
            raise PlanetNotFound()

        return PlanetDTO.from_model(planet)


@dataclass(frozen=True)
class GetSystemPlanets:
    repo: PlanetRepo
    system_repo: SystemRepo
    idp: IdentityProvider

    async def execute(self, data: str) -> list[PlanetDTO]:
        await self.idp.get_current_user()
        system = await self.system_repo.find_one_system(data)
        if not system:
            raise SystemNotFound()

        planets = await self.repo.find_all_planets_by_system(system.id)

        return [PlanetDTO.from_model(p) for p in planets]
