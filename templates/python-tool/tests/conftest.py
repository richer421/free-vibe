import pytest
from unittest.mock import AsyncMock, MagicMock

from __MODULE_NAME__.tools.registry import ToolContext


@pytest.fixture
def mock_ctx() -> ToolContext:
    """A ToolContext with mocked db and redis for unit tests."""
    return ToolContext(
        db=AsyncMock(),
        redis=AsyncMock(),
    )
