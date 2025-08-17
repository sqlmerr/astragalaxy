from dataclasses import dataclass
from uuid import uuid4

from voidspace.database.models import System
from voidspace.database.models.system_connection import SystemConnection
from voidspace.dto.system import CreateSystemDTO, SystemDTO
from voidspace.exceptions import AccessDeniedError
from voidspace.exceptions.system import SystemNotFound
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.interfaces.system_connection.repo import SystemConnectionRepo
from voidspace.use_cases import BaseUseCase
from voidspace.utils import generate_random_id


@dataclass(frozen=True)
class CreateSystem(BaseUseCase[CreateSystemDTO, SystemDTO]):
    repo: SystemRepo
    connection_repo: SystemConnectionRepo
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
        for conn in data.connections:
            sys = await self.repo.find_one_system(conn)
            if not sys:
                raise SystemNotFound()
            self.connection_repo.add_connection(
                SystemConnection(
                    id=uuid4(), system_from_id=system.id, system_to_id=sys.id
                )
            )
            self.connection_repo.add_connection(
                SystemConnection(
                    id=uuid4(), system_from_id=sys.id, system_to_id=system.id
                )
            )

        dto = SystemDTO(
            id=system.id,
            connections=data.connections,
            name=system.name,
            locations=system.locations,
        )

        return dto
