from typing import Protocol

from astragalaxy.database.models import System


class SystemRepo(Protocol):
    def add(self, system: System) -> None:
        raise NotImplementedError

    async def find_one_system(self, id: str) -> System | None:
        raise NotImplementedError

    async def find_one_random_system(self) -> System | None:
        raise NotImplementedError

    async def find_one_system_by_name(self, name: str) -> System | None:
        raise NotImplementedError

    async def find_all_systems(self, limit: int, offset: int) -> list[System]:
        raise NotImplementedError

    async def delete_system(self, id: str) -> None:
        raise NotImplementedError
