from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from voidspace.api.dependencies import JwtSecurity, CharacterSecurity
from voidspace.api.schemas.cooldown import CooldownSchema
from voidspace.api.schemas.navigation import PlanetNavigationSchema, HyperjumpSchema
from voidspace.use_cases.hyperjump import Hyperjump
from voidspace.use_cases.navigate_to_planet import NavigateToPlanet

router = APIRouter(prefix="/navigation", route_class=DishkaRoute, tags=["Navigation"])


@router.post("/planet", dependencies=[JwtSecurity, CharacterSecurity])
async def navigate_to_planet(
    data: PlanetNavigationSchema, use_case: FromDishka[NavigateToPlanet]
) -> CooldownSchema:
    cooldown = await use_case.execute(data.planet_id)
    return CooldownSchema.from_dto(cooldown)


@router.post("/hyperjump", dependencies=[JwtSecurity, CharacterSecurity])
async def hyperjump(
    data: HyperjumpSchema, use_case: FromDishka[Hyperjump]
) -> CooldownSchema:
    cooldown = await use_case.execute(data.path)
    return CooldownSchema.from_dto(cooldown)
