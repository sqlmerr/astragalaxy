from fastapi import FastAPI

from voidspace.config import Settings
from voidspace.database import init_db
from voidspace.di import init_di
from voidspace.routes import v1_router

from dishka.integrations.fastapi import setup_dishka


def create_app() -> FastAPI:
    config = Settings()
    app = FastAPI(title="VoidSpace")

    @app.get("/")
    async def root():
        return {"ok": True, "message": "Everything is fine"}

    app.include_router(v1_router)

    session_maker = init_db(config)
    container = init_di(config, session_maker)
    setup_dishka(container, app)

    return app
