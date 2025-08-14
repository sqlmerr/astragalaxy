from dishka.integrations.fastapi import setup_dishka
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from voidspace.api.routes import v1_router
from voidspace.config import Settings
from voidspace.database import init_db
from voidspace.di import init_di
from voidspace.exceptions import AppError

config = Settings()
app = FastAPI(title="VoidSpace")


@app.get("/")
async def root():
    return {"ok": True, "message": "Everything is fine"}


app.include_router(v1_router)

session_maker = init_db(config)
container = init_di(config, session_maker)
setup_dishka(container, app)


@app.exception_handler(AppError)
async def error_handler(request: Request, exc: AppError):
    return JSONResponse(status_code=exc.status, content={"message": exc.message})
