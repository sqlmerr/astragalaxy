from dataclasses import dataclass
import os

from dotenv import load_dotenv


@dataclass(frozen=True, kw_only=True)
class Settings:
    environment: str = "prod"  # prod, dev, test

    postgres_host: str
    postgres_port: int
    postgres_user: str
    postgres_password: str
    postgres_database: str = ""

    redis_url: str

    jwt_secret: str

    def build_postgres_dsn(self) -> str:
        return f"postgresql+asyncpg://{self.postgres_user}:{self.postgres_password}@{self.postgres_host}:{self.postgres_port}/{self.postgres_database}"


def load_settings_from_env() -> Settings:
    env = os.getenv("ENVIRONMENT")
    if not env:
        env = "prod"
    if env == "prod":
        load_dotenv(".env")
    elif env == "dev":
        load_dotenv(".env.dev")
    elif env == "test":
        load_dotenv(".env.test")
    else:
        raise ValueError(f"Invalid environment: {env}. Must be one of: prod, dev, test")

    return Settings(
        environment=env,
        postgres_host=os.environ["POSTGRES_HOST"],
        postgres_port=int(os.environ["POSTGRES_PORT"]),
        postgres_user=os.environ["POSTGRES_USER"],
        postgres_password=os.environ["POSTGRES_PASSWORD"],
        postgres_database=os.environ["POSTGRES_DATABASE"],
        redis_url=os.environ["REDIS_URL"],
        jwt_secret=os.environ["JWT_SECRET"],
    )
