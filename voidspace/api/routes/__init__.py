from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from . import auth, characters

v1_router = APIRouter(prefix="/v1", route_class=DishkaRoute)

v1_router.include_router(auth.router)
v1_router.include_router(characters.router)
