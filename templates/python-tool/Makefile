.PHONY: install dev-http dev-mcp dev-cli migrate test lint

install:
	uv sync

dev-http:
	uv run __MODULE_NAME__-http

dev-mcp:
	uv run __MODULE_NAME__-mcp

dev-cli:
	uv run __MODULE_NAME__ --help

migrate:
	uv run alembic upgrade head

migrate-new:
	uv run alembic revision --autogenerate -m "$(msg)"

test:
	uv run pytest

lint:
	uv run ruff check src tests
