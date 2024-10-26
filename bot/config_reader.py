from typing import Type

from pydantic import BaseModel
from pydantic_settings import (
    BaseSettings,
    SettingsConfigDict,
    PydanticBaseSettingsSource,
    TomlConfigSettingsSource,
)


class Settings(BaseSettings):
    bot_token: str
    api_url: str

    model_config = SettingsConfigDict(env_ignore_empty=True)

config = Settings()
