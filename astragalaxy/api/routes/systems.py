from typing import Annotated

from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from astragalaxy.api.dependencies import JwtSecurity, PaginationDepends
from astragalaxy.api.schemas import Pagination, DataSchema
from astragalaxy.api.schemas.planet import PlanetSchema
from astragalaxy.api.schemas.system import SystemSchema
from astragalaxy.dto.common import PaginationDTO
from astragalaxy.use_cases.get_planet import GetSystemPlanets
from astragalaxy.use_cases.get_system import GetSystemsPaginated, GetSystem

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


@router.get("/{id}/planets", dependencies=[JwtSecurity])
async def get_system_planets(
    id: str, use_case: FromDishka[GetSystemPlanets]
) -> DataSchema[PlanetSchema]:
    planets = await use_case.execute(id)

    return DataSchema(data=[PlanetSchema.from_dto(p) for p in planets])
