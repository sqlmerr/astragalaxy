from pydantic_settings import (
    BaseSettings,
    SettingsConfigDict,
)


class Settings(BaseSettings):
    bot_token: str
    api_url: str
    secret_token: str
    redis_url: str

    admins: tuple[int] = (1341947575,)

    model_config = SettingsConfigDict(env_ignore_empty=True, env_file=".env")


config = Settings()
