from typing import Protocol
from uuid import UUID

from voidspace.database.models import Inventory
from voidspace.database.models.inventory import InventoryOwnerEnum


class InventoryRepo(Protocol):
    def add_inventory(self, inventory: Inventory) -> None:
        raise NotImplementedError

    async def find_one_inventory(self, id: UUID) -> Inventory | None:
        raise NotImplementedError

    async def find_one_inventory_by_owner(
        self, owner_type: InventoryOwnerEnum, owner_id: UUID
    ) -> Inventory | None:
        raise NotImplementedError
