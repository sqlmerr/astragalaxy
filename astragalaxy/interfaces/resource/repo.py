from typing import Protocol
from uuid import UUID

from astragalaxy.database.models import Resource


class ResourceRepo(Protocol):
    def add_resource(self, resource: Resource) -> None:
        raise NotImplementedError

    async def find_one_resource(self, id: UUID) -> Resource | None:
        raise NotImplementedError

    async def find_all_resources_by_inventory_id(
        self, inventory_id: UUID
    ) -> list[Resource]:
        raise NotImplementedError

    async def find_one_resource_by_code_and_inventory(
        self, code: str, inventory_id: UUID
    ) -> Resource | None:
        raise NotImplementedError

    async def delete_resource(self, id: UUID) -> None:
        raise NotImplementedError
