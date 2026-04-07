from fastmcp import FastMCP

from ..storage.db import get_db_session
from ..storage.redis import get_redis_client
from ..tools.registry import ToolContext, ToolDef, discover, get_registry

_TOOLS_PACKAGE = "__MODULE_NAME__.tools"


async def _build_context() -> ToolContext:
    db_gen = get_db_session()
    redis_gen = get_redis_client()
    db = await db_gen.__anext__()
    redis = await redis_gen.__anext__()
    return ToolContext(db=db, redis=redis)


def _wrap(tool_def: ToolDef):
    """Wrap a tool function so fastmcp sees the correct Pydantic schema."""
    input_model = tool_def.input_model
    fn = tool_def.fn

    async def wrapped(**kwargs):
        params = input_model(**kwargs)
        ctx = await _build_context()
        result = await fn(params, ctx)
        return result.model_dump()

    wrapped.__name__ = tool_def.name
    wrapped.__doc__ = tool_def.description
    # Expose individual fields so fastmcp can generate a proper schema.
    wrapped.__annotations__ = {
        k: v.annotation
        for k, v in input_model.model_fields.items()
    }
    return wrapped


def create_server() -> FastMCP:
    discover(_TOOLS_PACKAGE)

    mcp = FastMCP("__MODULE_NAME__")

    for tool_def in get_registry().values():
        mcp.add_tool(_wrap(tool_def), name=tool_def.name, description=tool_def.description)

    return mcp
