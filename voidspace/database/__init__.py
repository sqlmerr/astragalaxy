from typing import Any

from redis.asyncio import Redis
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession
from sqlalchemy.orm import DeclarativeBase

from voidspace.config import Settings


class Base(DeclarativeBase):
    type_annotation_map = {dict[str, Any]: JSONB}


def init_db(config: Settings) -> async_sessionmaker[AsyncSession]:
    engine = create_async_engine(config.build_postgres_dsn())

    return async_sessionmaker(bind=engine)


def init_redis(config: Settings) -> Redis:
    return Redis.from_url(config.REDIS_URL)
