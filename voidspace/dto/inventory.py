from dataclasses import dataclass
from uuid import UUID

from voidspace.database.models.inventory import InventoryOwnerEnum, Inventory


@dataclass(frozen=True)
class InventoryDTO:
    id: UUID
    owner: InventoryOwnerEnum
    owner_id: UUID

    @classmethod
    def from_model(cls, model: Inventory) -> "InventoryDTO":
        return cls(id=model.id, owner_id=model.owner_id, owner=model.owner)
