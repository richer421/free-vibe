from pydantic import BaseModel
from pydantic_settings import (
    BaseSettings,
    PydanticBaseSettingsSource,
    SettingsConfigDict,
    YamlConfigSettingsSource,
)


class DatabaseConfig(BaseModel):
    url: str = "mysql+aiomysql://root:password@localhost:3306/__MODULE_NAME__"


class HttpConfig(BaseModel):
    host: str = "0.0.0.0"
    port: int = 8000


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_nested_delimiter="__")

    database: DatabaseConfig = DatabaseConfig()
    http: HttpConfig = HttpConfig()

    @classmethod
    def settings_customise_sources(
        cls,
        settings_cls: type[BaseSettings],
        env_settings: PydanticBaseSettingsSource,
        **kwargs,
    ) -> tuple[PydanticBaseSettingsSource, ...]:
        return (env_settings, YamlConfigSettingsSource(settings_cls, yaml_file="config.yaml"))


settings = Settings()
