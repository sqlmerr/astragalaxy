from dataclasses import dataclass

from astragalaxy.dto.common import PaginationDTO
from astragalaxy.dto.system import SystemDTO
from astragalaxy.exceptions.system import SystemNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.system.repo import SystemRepo
from astragalaxy.interfaces.system_connection.repo import SystemConnectionRepo


@dataclass(frozen=True)
class GetSystem:
    repo: SystemRepo
    idp: IdentityProvider
    connection_repo: SystemConnectionRepo

    async def execute(self, data: str) -> SystemDTO:
        await self.idp.get_current_user()

        s = await self.repo.find_one_system(data)
        if not s:
            raise SystemNotFound()

        conns = await self.connection_repo.find_all_connections_by_system_id(s.id)
        return SystemDTO(
            id=s.id,
            name=s.name,
            connections=[c.system_from_id for c in conns],
        )


@dataclass(frozen=True)
class GetSystemsPaginated:
    repo: SystemRepo
    connection_repo: SystemConnectionRepo
    idp: IdentityProvider

    async def execute(self, data: PaginationDTO) -> list[SystemDTO]:
        await self.idp.get_current_user()
        offset = 0 if data.page == 0 else data.page * data.per_page
        systems = await self.repo.find_all_systems(data.per_page, offset)
        dtos = []
        for s in systems:
            conns = await self.connection_repo.find_all_connections_by_system_id(s.id)
            dtos.append(
                SystemDTO(
                    id=s.id,
                    name=s.name,
                    connections=[c.system_from_id for c in conns],
                )
            )

        return dtos
