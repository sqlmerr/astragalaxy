from dataclasses import dataclass
from uuid import UUID

from sqlalchemy import select, and_
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models import Inventory
from astragalaxy.database.models.inventory import InventoryOwnerEnum
from astragalaxy.interfaces.inventory.repo import InventoryRepo


@dataclass(frozen=True)
class InventoryRepository(InventoryRepo):
    session: AsyncSession

    def add_inventory(self, inventory: Inventory) -> None:
        self.session.add(inventory)

    async def find_one_inventory(self, id: UUID) -> Inventory | None:
        stmt = select(Inventory).where(Inventory.id == id)
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()

    async def find_one_inventory_by_owner(
        self, owner_type: InventoryOwnerEnum, owner_id: UUID
    ) -> Inventory | None:
        stmt = select(Inventory).where(
            and_(Inventory.owner == owner_type, Inventory.owner_id == owner_id)
        )
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()
