from collections.abc import AsyncGenerator

import pytest
from httpx import ASGITransport, AsyncClient

from __MODULE_NAME__.main import create_app
from __MODULE_NAME__.storage.db import close_engine, init_engine


@pytest.fixture
async def client() -> AsyncGenerator[AsyncClient, None]:
    init_engine()
    app = create_app()
    async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as c:
        yield c
    await close_engine()
