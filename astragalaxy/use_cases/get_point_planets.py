from dataclasses import dataclass

from astragalaxy.dto.planet import PlanetDTO
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.planet.repo import PlanetRepo


@dataclass(frozen=True)
class GetPointPlanets:
    repo: PlanetRepo
    idp: IdentityProvider

    async def execute(self, point_id: str) -> list[PlanetDTO]:
        await self.idp.get_current_user()

        planets = await self.repo.get_planets_by_point(point_id)

        return [PlanetDTO.from_model(p) for p in planets]
