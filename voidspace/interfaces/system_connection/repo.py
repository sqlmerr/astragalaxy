from typing import Protocol
from uuid import UUID

from voidspace.database.models.system_connection import SystemConnection


class SystemConnectionRepo(Protocol):
    def add_connection(self, data: SystemConnection) -> None:
        raise NotImplementedError

    async def find_one_connection(self, id: UUID) -> SystemConnection | None:
        raise NotImplementedError

    async def find_all_connections_by_system_id(
        self, system_id: str
    ) -> list[SystemConnection]:
        raise NotImplementedError

    async def find_connections_by_system_ids(
        self, system_a_id: str, system_b_id: str
    ) -> list[SystemConnection]:
        raise NotImplementedError

    async def delete_connection(self, id: UUID) -> None:
        raise NotImplementedError
