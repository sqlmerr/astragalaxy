from dataclasses import dataclass

from voidspace.database.models import System
from voidspace.dto.system import CreateSystemDTO, SystemDTO
from voidspace.exceptions import AccessDeniedError
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.use_cases import BaseUseCase
from voidspace.utils import generate_random_id


@dataclass(frozen=True)
class CreateSystem(BaseUseCase[CreateSystemDTO, SystemDTO]):
    repo: SystemRepo
    idp: IdentityProvider

    async def execute(self, data: CreateSystemDTO) -> SystemDTO:
        current_user = await self.idp.get_current_user()
        if current_user.username != "admin":  # TODO: roles
            raise AccessDeniedError()

        system_id = generate_random_id(8)
        await self.repo.create_system(
            System(id=system_id, name=data.name, locations=data.locations)
        )

        system = await self.repo.find_one_system(system_id)
        return SystemDTO.from_model(system)
