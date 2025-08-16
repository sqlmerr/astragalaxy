from dataclasses import dataclass

from voidspace.database.models import Planet
from voidspace.dto.planet import CreatePlanetDTO, PlanetDTO
from voidspace.exceptions import AccessDeniedError
from voidspace.exceptions.system import SystemNotFound
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.planet.repo import PlanetRepo
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.use_cases import BaseUseCase
from voidspace.utils import generate_random_id


@dataclass(frozen=True)
class CreatePlanet(BaseUseCase[CreatePlanetDTO, PlanetDTO]):
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
        return PlanetDTO.from_model(p)
