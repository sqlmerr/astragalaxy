from dataclasses import dataclass

from voidspace.dto.planet import PlanetDTO
from voidspace.exceptions.planet import PlanetNotFound
from voidspace.exceptions.system import SystemNotFound
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.planet.repo import PlanetRepo
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.use_cases import BaseUseCase


@dataclass(frozen=True)
class GetPlanet(BaseUseCase[str, PlanetDTO]):
    repo: PlanetRepo
    idp: IdentityProvider

    async def execute(self, data: str) -> PlanetDTO:
        await self.idp.get_current_user()

        planet = await self.repo.find_one_planet(data)
        if not planet:
            raise PlanetNotFound()

        return PlanetDTO.from_model(planet)


@dataclass(frozen=True)
class GetSystemPlanets(BaseUseCase[str, list[PlanetDTO]]):
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
