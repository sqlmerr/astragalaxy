from dataclasses import dataclass
from uuid import uuid4

from astragalaxy.database.models import System
from astragalaxy.database.models.system_connection import SystemConnection
from astragalaxy.dto.system import CreateSystemDTO, SystemDTO
from astragalaxy.exceptions import AccessDeniedError
from astragalaxy.exceptions.system import SystemNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.system.repo import SystemRepo
from astragalaxy.interfaces.system_connection.repo import SystemConnectionRepo
from astragalaxy.utils import generate_random_id


@dataclass(frozen=True)
class CreateSystem:
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
        if not system:
            raise SystemNotFound()
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
