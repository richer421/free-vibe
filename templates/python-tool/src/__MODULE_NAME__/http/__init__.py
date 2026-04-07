from .app import create_app


def main() -> None:
    import uvicorn
    from ..config import settings

    uvicorn.run(create_app(), host=settings.http_host, port=settings.http_port)
