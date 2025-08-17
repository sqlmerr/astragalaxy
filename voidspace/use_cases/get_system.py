from dataclasses import dataclass

from voidspace.dto.common import PaginationDTO
from voidspace.dto.system import SystemDTO
from voidspace.exceptions.system import SystemNotFound
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.interfaces.system_connection.repo import SystemConnectionRepo
from voidspace.use_cases import BaseUseCase


@dataclass(frozen=True)
class GetSystem(BaseUseCase[str, SystemDTO]):
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
            locations=s.locations,
            connections=[c.system_from_id for c in conns],
        )


@dataclass(frozen=True)
class GetSystemsPaginated(BaseUseCase[PaginationDTO, list[SystemDTO]]):
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
                    locations=s.locations,
                    connections=[c.system_from_id for c in conns],
                )
            )

        return dtos
