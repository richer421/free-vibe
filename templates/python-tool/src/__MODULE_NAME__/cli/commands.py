import asyncio
import json
import sys

import typer

from ..storage.db import get_db_session
from ..storage.redis import get_redis_client
from ..tools.registry import ToolContext, ToolDef, discover, get_registry

_TOOLS_PACKAGE = "__MODULE_NAME__.tools"


async def _run_tool(tool_def: ToolDef, params_json: str) -> str:
    params = tool_def.input_model.model_validate_json(params_json)
    db_gen = get_db_session()
    redis_gen = get_redis_client()
    db = await db_gen.__anext__()
    redis = await redis_gen.__anext__()
    ctx = ToolContext(db=db, redis=redis)
    result = await tool_def.fn(params, ctx)
    return result.model_dump_json(indent=2)


def _make_command(tool_def: ToolDef):
    def command(
        params: str = typer.Argument(
            ..., help=f"JSON input matching schema: {tool_def.input_model.model_json_schema()}"
        ),
    ) -> None:
        output = asyncio.run(_run_tool(tool_def, params))
        typer.echo(output)

    command.__name__ = tool_def.name
    command.__doc__ = tool_def.description
    return command


def create_app() -> typer.Typer:
    discover(_TOOLS_PACKAGE)

    app = typer.Typer(help="__MODULE_NAME__ tool runner")

    for tool_def in get_registry().values():
        app.command(name=tool_def.name, help=tool_def.description)(_make_command(tool_def))

    return app
