from contextlib import asynccontextmanager

from dishka.integrations.fastapi import setup_dishka
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from astragalaxy.api.routes import v1_router
from astragalaxy.config import load_settings_from_env
from astragalaxy.database import init_db, init_redis
from astragalaxy.database.models import System
from astragalaxy.di import init_di
from astragalaxy.exceptions import AppError
from astragalaxy.interfaces.system.repo import SystemRepo
from astragalaxy.utils import generate_random_id

config = load_settings_from_env()
session_maker = init_db(config)
redis = init_redis(config)
container = init_di(config, session_maker, redis)


@asynccontextmanager
async def lifespan(app: FastAPI):
    async with container() as nested_container:
        system_repo = await nested_container.get(SystemRepo)
        system = await system_repo.find_one_system_by_name("initial")
        if not system:
            await system_repo.create_system(
                System(
                    id=generate_random_id(8),
                    name="initial",
                    locations=[],
                )
            )
    yield


app = FastAPI(title="AstraGalaxy", lifespan=lifespan)
setup_dishka(container, app)


@app.get("/")
async def root():
    return {"ok": True, "message": "Everything is fine"}


app.include_router(v1_router)


@app.exception_handler(AppError)
async def error_handler(request: Request, exc: AppError):
    return JSONResponse(status_code=exc.status, content={"message": exc.message})
