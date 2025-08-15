from typing import Protocol

from voidspace.database.models import System


class SystemRepo(Protocol):
    async def create_system(self, system: System) -> str:
        raise NotImplementedError

    async def find_one_system(self, id: str) -> System | None:
        raise NotImplementedError

    async def find_all_systems(self, limit: int, offset: int) -> list[System]:
        raise NotImplementedError

    async def delete_system(self, id: str) -> None:
        raise NotImplementedError
