from typing import Annotated

from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from astragalaxy.api.dependencies import JwtSecurity, PaginationDepends
from astragalaxy.api.schemas import Pagination, DataSchema
from astragalaxy.api.schemas.planet import PlanetSchema
from astragalaxy.api.schemas.point import PointSchema
from astragalaxy.api.schemas.system import SystemSchema
from astragalaxy.dto.common import PaginationDTO
from astragalaxy.use_cases.get_system import GetSystemsPaginated, GetSystem
from astragalaxy.use_cases.get_system_points import GetSystemPoints

router = APIRouter(prefix="/systems", route_class=DishkaRoute, tags=["Systems"])


@router.get("/", dependencies=[JwtSecurity])
async def get_systems_paginated(
    data: Annotated[Pagination, PaginationDepends],
    use_case: FromDishka[GetSystemsPaginated],
) -> DataSchema[SystemSchema]:
    systems = await use_case.execute(
        PaginationDTO(per_page=data.per_page, page=data.page)
    )
    schemas = [SystemSchema.from_dto(s) for s in systems]

    return DataSchema(data=schemas)


@router.get("/{id}", dependencies=[JwtSecurity])
async def get_system(id: str, use_case: FromDishka[GetSystem]) -> SystemSchema:
    system = await use_case.execute(id)

    return SystemSchema.from_dto(system)


@router.get("/{id}/points", dependencies=[JwtSecurity])
async def get_system_points(id: str, use_case: FromDishka[GetSystemPoints]) -> DataSchema[PointSchema]:
    res = await use_case.execute(id)

    return DataSchema(data=[PointSchema.from_dto(p) for p in res])
