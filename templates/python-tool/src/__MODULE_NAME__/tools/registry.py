import importlib
import inspect
import pkgutil
from dataclasses import dataclass
from typing import Any, Callable, get_type_hints

from pydantic import BaseModel
from sqlalchemy.ext.asyncio import AsyncSession
from redis.asyncio import Redis


@dataclass
class ToolContext:
    """Injected into every tool function. Provides access to shared storage."""

    db: AsyncSession
    redis: Redis


@dataclass
class ToolDef:
    name: str
    description: str
    fn: Callable
    input_model: type[BaseModel]
    output_model: type[BaseModel]


_registry: dict[str, ToolDef] = {}


def tool(name: str | None = None, description: str | None = None) -> Callable:
    """Decorator that registers an async tool function into the global registry.

    The decorated function must follow the signature:
        async def my_tool(params: InputModel, ctx: ToolContext) -> OutputModel
    """

    def decorator(fn: Callable) -> Callable:
        tool_name = name or fn.__name__
        tool_desc = (description or inspect.getdoc(fn) or "").strip()

        hints = get_type_hints(fn)
        sig = inspect.signature(fn)
        params_list = list(sig.parameters.keys())

        input_model: type[BaseModel] | None = None
        for param_name in params_list:
            if param_name == "ctx":
                continue
            hint = hints.get(param_name)
            if hint and isinstance(hint, type) and issubclass(hint, BaseModel):
                input_model = hint
                break

        output_model = hints.get("return")

        if input_model is None:
            raise TypeError(f"Tool '{tool_name}' must have a Pydantic BaseModel as first parameter")
        if output_model is None or not issubclass(output_model, BaseModel):
            raise TypeError(f"Tool '{tool_name}' must return a Pydantic BaseModel")

        _registry[tool_name] = ToolDef(
            name=tool_name,
            description=tool_desc,
            fn=fn,
            input_model=input_model,
            output_model=output_model,
        )
        return fn

    return decorator


def get_registry() -> dict[str, ToolDef]:
    return dict(_registry)


def discover(package: str) -> None:
    """Import all modules in the given package to trigger @tool registrations."""
    mod = importlib.import_module(package)
    for _, module_name, _ in pkgutil.iter_modules(mod.__path__):  # type: ignore[union-attr]
        if module_name != "registry":
            importlib.import_module(f"{package}.{module_name}")
