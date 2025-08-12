from fastapi import FastAPI

from voidspace.routes import v1_router


def create_app() -> FastAPI:
    app = FastAPI(title="VoidSpace")

    @app.get("/")
    async def root():
        return {"ok": True, "message": "Everything is fine"}

    app.include_router(v1_router)

    return app
