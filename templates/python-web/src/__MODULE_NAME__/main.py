from contextlib import asynccontextmanager

from fastapi import FastAPI

from .api.router import router
from .storage.db import close_engine, init_engine


@asynccontextmanager
async def lifespan(app: FastAPI):
    init_engine()
    yield
    await close_engine()


def create_app() -> FastAPI:
    app = FastAPI(title="__MODULE_NAME__", version="0.1.0", lifespan=lifespan)
    app.include_router(router, prefix="/api/v1")
    return app


def main() -> None:
    import uvicorn
    from .config import settings

    uvicorn.run(create_app(), host=settings.http.host, port=settings.http.port)
