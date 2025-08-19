from typing import Protocol
from uuid import UUID

from voidspace.database.models import Item


class ItemRepo(Protocol):
    def add_item(self, item: Item) -> None:
        raise NotImplementedError

    async def find_one_item(self, id: UUID) -> Item | None:
        raise NotImplementedError

    async def find_all_items_by_inventory_id(self, inventory_id: UUID) -> list[Item]:
        raise NotImplementedError

    async def find_all_items_by_inventory_id_and_code(
        self, inventory_id: UUID, code: str
    ) -> list[Item]:
        raise NotImplementedError

    async def delete_item(self, id: UUID) -> None:
        raise NotImplementedError
