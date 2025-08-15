from typing import Protocol

from voidspace.dto.system import CreateSystemDTO, SystemDTO


class SystemWriter(Protocol):
    async def create_system(self, data: CreateSystemDTO) -> SystemDTO:
        raise NotImplementedError
