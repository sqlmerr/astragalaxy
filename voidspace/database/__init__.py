from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession
from sqlalchemy.orm import DeclarativeBase

from voidspace.config import Settings

class Base(DeclarativeBase):
    pass


async def init_db(config: Settings) -> async_sessionmaker[AsyncSession]:
    engine = create_async_engine(config.build_postgres_dsn())

    return async_sessionmaker(bind=engine)
