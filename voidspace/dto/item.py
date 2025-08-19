from dataclasses import dataclass
from typing import Any
from uuid import UUID

from voidspace.database.models import Item


@dataclass(frozen=True)
class ItemDTO:
    id: UUID
    code: str
    inventory_id: UUID
    data: dict[str, Any]
    durabilty: int

    @classmethod
    def from_model(cls, model: Item) -> "ItemDTO":
        return cls(
            id=model.id,
            code=model.code,
            inventory_id=model.inventory_id,
            data=model.data,
            durabilty=model.durability,
        )


@dataclass(frozen=True)
class TransferItemDTO:
    item_id: UUID
    receiver_inventory_id: UUID
