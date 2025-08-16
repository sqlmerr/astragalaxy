from dataclasses import dataclass

from voidspace.dto.common import PaginationDTO
from voidspace.dto.system import SystemDTO
from voidspace.exceptions.system import SystemNotFound
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.use_cases import BaseUseCase


@dataclass(frozen=True)
class GetSystem(BaseUseCase[str, SystemDTO]):
    repo: SystemRepo
    idp: IdentityProvider

    async def execute(self, data: str) -> SystemDTO:
        await self.idp.get_current_user()

        system = await self.repo.find_one_system(data)
        if not system:
            raise SystemNotFound()

        return SystemDTO.from_model(system)


@dataclass(frozen=True)
class GetSystemsPaginated(BaseUseCase[PaginationDTO, list[SystemDTO]]):
    repo: SystemRepo
    idp: IdentityProvider

    async def execute(self, data: PaginationDTO) -> list[SystemDTO]:
        await self.idp.get_current_user()
        offset = 0 if data.page == 0 else data.page * data.per_page
        systems = await self.repo.find_all_systems(data.per_page, offset)
        return [SystemDTO.from_model(s) for s in systems]
