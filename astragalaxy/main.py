from contextlib import asynccontextmanager

from dishka.integrations.fastapi import setup_dishka
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from astragalaxy.api.routes import v1_router
from astragalaxy.config import load_settings_from_env
from astragalaxy.database import init_db, init_redis
from astragalaxy.database.models import Point, Station, System
from astragalaxy.di import init_di
from astragalaxy.exceptions import AppError
from astragalaxy.interfaces.point.repo import PointRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.station.repo import StationRepo
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
            point_repo = await nested_container.get(PointRepo)
            station_repo = await nested_container.get(StationRepo)
            commiter = await nested_container.get(Commiter)

            system_id = generate_random_id(8)
            point_id = generate_random_id(16)

            system_repo.add(
                System(
                    id=system_id,
                    name="initial",
                )
            )
            point_repo.add(Point(
                id=point_id,
                name="station_point",
                system_id=system_id
            ))
            station_repo.add(Station(
                id=generate_random_id(8),
                point_id=point_id
            ))
            await commiter.commit()
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
