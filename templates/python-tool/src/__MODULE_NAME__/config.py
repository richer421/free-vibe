from pydantic import BaseModel
from pydantic_settings import BaseSettings, YamlConfigSettingsSource, PydanticBaseSettingsSource


class DatabaseConfig(BaseModel):
    url: str = "mysql+aiomysql://root:password@localhost:3306/__MODULE_NAME__"


class RedisConfig(BaseModel):
    url: str = "redis://localhost:6379/0"


class HttpConfig(BaseModel):
    host: str = "0.0.0.0"
    port: int = 8000


class Settings(BaseSettings):
    database: DatabaseConfig = DatabaseConfig()
    redis: RedisConfig = RedisConfig()
    http: HttpConfig = HttpConfig()

    @classmethod
    def settings_customise_sources(
        cls,
        settings_cls: type[BaseSettings],
        **kwargs,
    ) -> tuple[PydanticBaseSettingsSource, ...]:
        return (YamlConfigSettingsSource(settings_cls, yaml_file="config.yaml"),)


settings = Settings()
