from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from . import auth, characters, systems, spaceships, navigation, inventories, points, planets

v1_router = APIRouter(prefix="/v1", route_class=DishkaRoute)

v1_router.include_router(auth.router)
v1_router.include_router(characters.router)
v1_router.include_router(systems.router)
v1_router.include_router(planets.router)
v1_router.include_router(spaceships.router)
v1_router.include_router(navigation.router)
v1_router.include_router(inventories.router)
v1_router.include_router(points.router)