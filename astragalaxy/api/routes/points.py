from dishka.integrations.fastapi import DishkaRoute, FromDishka
from fastapi import APIRouter

from astragalaxy.api.dependencies import JwtSecurity
from astragalaxy.api.schemas import DataSchema
from astragalaxy.api.schemas.planet import PlanetSchema
from astragalaxy.api.schemas.point import CreatePointSchema, PointSchema
from astragalaxy.api.schemas.station import StationSchema
from astragalaxy.use_cases.create_point import CreatePoint
from astragalaxy.use_cases.get_point import GetPoint
from astragalaxy.use_cases.get_point_planets import GetPointPlanets
from astragalaxy.use_cases.get_point_stations import GetPointStations


router = APIRouter(prefix="/points", route_class=DishkaRoute, tags=["Points"])


@router.get("/{point_id}", dependencies=[JwtSecurity])
async def get_point(point_id: str, use_case: FromDishka[GetPoint]) -> PointSchema:
    res = await use_case.execute(point_id)
    return PointSchema.from_dto(res)


@router.get("/{id}/planets", dependencies=[JwtSecurity])
async def get_point_planets(point_id: str, use_case: FromDishka[GetPointPlanets]) -> DataSchema[PlanetSchema]:
    res = await use_case.execute(point_id)
    return DataSchema(data=[PlanetSchema.from_dto(p) for p in res])

@router.get("/{id}/stations", dependencies=[JwtSecurity])
async def get_point_stations(point_id: str, use_case: FromDishka[GetPointStations]) -> DataSchema[StationSchema]:
    res = await use_case.execute(point_id)
    return DataSchema(data=[StationSchema.from_dto(s) for s in res])


@router.post("/", dependencies=[JwtSecurity])
async def create_point(data: CreatePointSchema, use_case: FromDishka[CreatePoint]) -> PointSchema:
    """Only admins can create new points"""
    
    res = await use_case.execute(data.into_dto())
    return PointSchema.from_dto(res)
