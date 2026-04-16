from fastapi import APIRouter

from .v1 import health

router = APIRouter()
router.include_router(health.router, tags=["health"])
