from dataclasses import dataclass

from astragalaxy.database.models import Planet
from astragalaxy.dto.planet import CreatePlanetDTO, PlanetDTO
from astragalaxy.exceptions import AccessDeniedError, AppError
from astragalaxy.exceptions.point import PointNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.planet.repo import PlanetRepo
from astragalaxy.interfaces.point.repo import PointRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.system.repo import SystemRepo
from astragalaxy.utils import generate_random_id


@dataclass(frozen=True)
class CreatePlanet:
    repo: PlanetRepo
    point_repo: PointRepo
    idp: IdentityProvider
    commiter: Commiter

    async def execute(self, data: CreatePlanetDTO) -> PlanetDTO:
        current_user = await self.idp.get_current_user()
        if current_user.username != "admin":  # TODO: roles
            raise AccessDeniedError()

        point = await self.point_repo.find_one_point(data.point_id)
        if not point:
            raise PointNotFound()

        planet = Planet(
            id=generate_random_id(16),
            name=data.name,
            point_id=data.point_id,
            threat=data.threat,
        )
        self.repo.add(planet)
        dto = PlanetDTO.from_model(planet)
        await self.commiter.commit()

        return dto
