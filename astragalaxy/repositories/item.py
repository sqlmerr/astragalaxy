from dataclasses import dataclass
from uuid import UUID

from sqlalchemy import select, and_, delete
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models import Item
from astragalaxy.interfaces.item.repo import ItemRepo


@dataclass(frozen=True)
class ItemRepository(ItemRepo):
    session: AsyncSession

    def add_item(self, item: Item) -> None:
        self.session.add(item)

    async def find_one_item(self, id: UUID) -> Item | None:
        stmt = select(Item).where(Item.id == id)
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()

    async def find_all_items_by_inventory_id(self, inventory_id: UUID) -> list[Item]:
        stmt = select(Item).where(Item.inventory_id == inventory_id)
        result = await self.session.execute(stmt)

        return list(result.scalars().all())

    async def find_all_items_by_inventory_id_and_code(
        self, inventory_id: UUID, code: str
    ) -> list[Item]:
        stmt = select(Item).where(
            and_(Item.inventory_id == inventory_id, Item.code == code)
        )
        result = await self.session.execute(stmt)
        return list(result.scalars().all())

    async def delete_item(self, id: UUID) -> None:
        stmt = delete(Item).where(Item.id == id)
        await self.session.execute(stmt)
