from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    db_url: str = "mysql+aiomysql://root:password@localhost:3306/__MODULE_NAME__"
    redis_url: str = "redis://localhost:6379/0"
    http_host: str = "0.0.0.0"
    http_port: int = 8000

    model_config = {"env_file": ".env", "env_file_encoding": "utf-8"}


settings = Settings()
