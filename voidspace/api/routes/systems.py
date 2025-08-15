from typing import Annotated

from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from voidspace.api.dependencies import JwtSecurity, PaginationDepends
from voidspace.api.schemas import Pagination, DataSchema
from voidspace.api.schemas.system import SystemSchema
from voidspace.interfaces.system.read import SystemReader

router = APIRouter(prefix="/systems", route_class=DishkaRoute, tags=["Systems"])


@router.get("/", dependencies=[JwtSecurity])
async def get_systems_paginated(
    data: Annotated[Pagination, PaginationDepends],
    system_reader: FromDishka[SystemReader],
) -> DataSchema[SystemSchema]:
    systems = await system_reader.get_systems_paginated(data.per_page, data.page)
    schemas = [SystemSchema.from_dto(s) for s in systems]

    return DataSchema(data=schemas)
