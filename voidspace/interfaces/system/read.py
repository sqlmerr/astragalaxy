from typing import Protocol

from voidspace.dto.system import SystemDTO


class SystemReader(Protocol):
    async def get_one_system(self, id: str) -> SystemDTO | None:
        raise NotImplementedError

    async def get_systems_paginated(
        self, per_page: int = 10, page: int = 0
    ) -> list[SystemDTO]:
        raise NotImplementedError
