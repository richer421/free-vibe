.PHONY: install dev migrate migrate-new test lint format infra-up infra-down

install:
	uv sync

dev:
	uv run __MODULE_NAME__

migrate:
	uv run alembic upgrade head

migrate-new:
	uv run alembic revision --autogenerate -m "$(msg)"

test:
	uv run pytest

lint:
	uv run ruff check src tests

format:
	uv run ruff format src tests
	uv run ruff check --fix src tests

infra-up:
	docker compose up -d

infra-down:
	docker compose down
