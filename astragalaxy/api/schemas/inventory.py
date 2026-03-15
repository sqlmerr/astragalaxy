from typing import Any
from uuid import UUID

from pydantic import BaseModel

from astragalaxy.dto.item import ItemDTO
from astragalaxy.dto.resource import ResourceDTO


class ItemSchema(BaseModel):
    id: UUID
    code: str
    inventory_id: UUID
    data: dict[str, Any]
    durability: int

    @classmethod
    def from_dto(cls, dto: ItemDTO) -> "ItemSchema":
        return cls.model_validate(dto, from_attributes=True)



class ResourceSchema(BaseModel):
    id: UUID
    code: str
    quantity: int
    inventory_id: UUID

    @classmethod
    def from_dto(cls, dto: ResourceDTO) -> "ResourceSchema":
        return cls.model_validate(dto, from_attributes=True)

