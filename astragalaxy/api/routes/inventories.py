from uuid import UUID

from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from astragalaxy.api.schemas import DataSchema
from astragalaxy.api.schemas.inventory import ItemSchema, ResourceSchema
from astragalaxy.database.models.inventory import InventoryOwnerEnum
from astragalaxy.use_cases.get_inventory import (
    GetInventoryItems,
    GetInventoryByOwnerRequest,
    GetInventoryByIdRequest,
    GetInventoryResources,
)

router = APIRouter(prefix="/inventory", route_class=DishkaRoute, tags=["Inventory"])


@router.get("/{owner}/{id}/items")
async def get_inventory_items(
    owner: InventoryOwnerEnum, id: UUID, use_case: FromDishka[GetInventoryItems]
) -> DataSchema[ItemSchema]:
    items = await use_case.execute(GetInventoryByOwnerRequest(owner=owner, owner_id=id))
    return DataSchema(data=[ItemSchema.from_dto(i) for i in items])


@router.get("/{inventory_id}/items")
async def get_inventory_items_by_inventory_id(
    inventory_id: UUID, use_case: FromDishka[GetInventoryItems]
) -> DataSchema[ItemSchema]:
    items = await use_case.execute(GetInventoryByIdRequest(inventory_id=inventory_id))
    return DataSchema(data=[ItemSchema.from_dto(i) for i in items])


@router.get("/{owner}/{id}/resources")
async def get_inventory_resources(
    owner: InventoryOwnerEnum, id: UUID, use_case: FromDishka[GetInventoryResources]
) -> DataSchema[ResourceSchema]:
    resources = await use_case.execute(
        GetInventoryByOwnerRequest(owner=owner, owner_id=id)
    )
    return DataSchema(data=[ResourceSchema.from_dto(i) for i in resources])


@router.get("/{inventory_id}/resources")
async def get_inventory_resources_by_inventory_id(
    inventory_id: UUID, use_case: FromDishka[GetInventoryResources]
) -> DataSchema[ResourceSchema]:
    resources = await use_case.execute(
        GetInventoryByIdRequest(inventory_id=inventory_id)
    )
    return DataSchema(data=[ResourceSchema.from_dto(i) for i in resources])
