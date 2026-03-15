from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from astragalaxy.api.dependencies import JwtSecurity
from astragalaxy.api.schemas.planet import CreatePlanetSchema, PlanetSchema
from astragalaxy.use_cases.create_planet import CreatePlanet

router = APIRouter(prefix="/planets", route_class=DishkaRoute, tags=["Planets"])


@router.post("/", dependencies=[JwtSecurity])
async def create_planet(data: CreatePlanetSchema, use_case: FromDishka[CreatePlanet]) -> PlanetSchema:
    """Only admins can create new planets"""
    
    res = await use_case.execute(data.into_dto())

    return PlanetSchema.from_dto(res)
