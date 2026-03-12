from dataclasses import dataclass

from astragalaxy.database.models import Planet
from astragalaxy.dto.planet import CreatePlanetDTO, PlanetDTO
from astragalaxy.exceptions import AccessDeniedError, AppError
from astragalaxy.exceptions.system import SystemNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.planet.repo import PlanetRepo
from astragalaxy.interfaces.system.repo import SystemRepo
from astragalaxy.utils import generate_random_id


@dataclass(frozen=True)
class CreatePlanet:
    repo: PlanetRepo
    system_repo: SystemRepo
    idp: IdentityProvider

    async def execute(self, data: CreatePlanetDTO) -> PlanetDTO:
        current_user = await self.idp.get_current_user()
        if current_user.username != "admin":  # TODO: roles
            raise AccessDeniedError()

        system = await self.system_repo.find_one_system(data.system_id)
        if not system:
            raise SystemNotFound()

        planet = Planet(
            id=generate_random_id(16),
            name=data.name,
            system_id=data.system_id,
            threat=data.threat,
        )
        await self.repo.create_planet(planet)

        p = await self.repo.find_one_planet(planet.id)
        if not p:
            raise AppError()
        return PlanetDTO.from_model(p)
