from dataclasses import dataclass
from uuid import UUID

from voidspace.database.models import Resource


@dataclass(frozen=True)
class ResourceDTO:
    id: UUID
    code: str
    quantity: int
    inventory_id: UUID

    @classmethod
    def from_model(cls, model: Resource) -> "ResourceDTO":
        return cls(
            id=model.id,
            code=model.code,
            quantity=model.quantity,
            inventory_id=model.inventory_id,
        )
