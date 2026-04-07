from fastapi import Depends, FastAPI

from ..storage.db import get_db_session
from ..storage.redis import get_redis_client
from ..tools.registry import ToolContext, ToolDef, discover, get_registry

_TOOLS_PACKAGE = "__MODULE_NAME__.tools"


def _make_endpoint(tool_def: ToolDef):
    """Factory that creates a typed FastAPI endpoint for the given tool."""
    input_model = tool_def.input_model
    output_model = tool_def.output_model
    fn = tool_def.fn

    async def endpoint(
        params: input_model,  # type: ignore[valid-type]
        db=Depends(get_db_session),
        redis=Depends(get_redis_client),
    ) -> output_model:  # type: ignore[valid-type]
        ctx = ToolContext(db=db, redis=redis)
        return await fn(params, ctx)

    endpoint.__name__ = tool_def.name
    endpoint.__annotations__["params"] = input_model
    endpoint.__annotations__["return"] = output_model
    return endpoint


def create_app() -> FastAPI:
    discover(_TOOLS_PACKAGE)

    app = FastAPI(title=__MODULE_NAME__, version="0.1.0")

    for tool_def in get_registry().values():
        app.add_api_route(
            f"/tools/{tool_def.name}",
            _make_endpoint(tool_def),
            methods=["POST"],
            response_model=tool_def.output_model,
            summary=tool_def.description,
            tags=["tools"],
        )

    return app
