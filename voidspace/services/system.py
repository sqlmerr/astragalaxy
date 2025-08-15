from dataclasses import dataclass

from voidspace.database.models import System
from voidspace.dto.system import SystemDTO, CreateSystemDTO
from voidspace.exceptions.system import SystemNotFound
from voidspace.interfaces.system.create import SystemWriter
from voidspace.interfaces.system.delete import SystemDeleter
from voidspace.interfaces.system.read import SystemReader
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.utils import generate_random_id


@dataclass(frozen=True)
class SystemService(SystemReader, SystemWriter, SystemDeleter):
    repo: SystemRepo

    async def get_one_system(self, id: str) -> SystemDTO | None:
        system = await self.repo.find_one_system(id)
        if not system:
            raise SystemNotFound

        return SystemDTO.from_model(system)

    async def get_systems_paginated(
        self, per_page: int = 10, page: int = 0
    ) -> list[SystemDTO]:
        offset = 0 if page == 0 else page * per_page
        systems = await self.repo.find_all_systems(per_page, offset)
        return [SystemDTO.from_model(s) for s in systems]

    async def create_system(self, data: CreateSystemDTO) -> SystemDTO:
        system_id = generate_random_id(8)
        await self.repo.create_system(
            System(id=system_id, name=data.name, locations=data.locations)
        )

        return await self.get_one_system(system_id)

    async def delete_system(self, id: str) -> None:
        await self.repo.delete_system(id)
