from collections.abc import AsyncGenerator

from redis.asyncio import Redis

_client: Redis | None = None


def _get_client() -> Redis:
    global _client
    if _client is None:
        from ..config import settings
        _client = Redis.from_url(settings.redis_url, decode_responses=True)
    return _client


async def get_redis_client() -> AsyncGenerator[Redis, None]:
    yield _get_client()
