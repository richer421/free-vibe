import pytest

from __MODULE_NAME__.tools.example import EchoInput, echo
from __MODULE_NAME__.tools.registry import ToolContext


async def test_echo(mock_ctx: ToolContext) -> None:
    result = await echo(EchoInput(message="hello"), mock_ctx)
    assert result.echo == "hello"
