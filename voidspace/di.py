from typing import AsyncGenerator

from dishka import Provider, make_async_container, provide, Scope, AsyncContainer
from sqlalchemy.ext.asyncio import async_sessionmaker, AsyncSession
from sqlalchemy.exc import SQLAlchemyError

from voidspace.config import Settings


class ConfigProvider(Provider):
    def __init__(self, settings: Settings):
        super().__init__()
        self.settings = settings

    @provide(scope=Scope.APP)
    def get_settings(self) -> Settings:
        return self.settings

class DatabaseProvider(Provider):
    def __init__(self, session_maker: async_sessionmaker):
        super().__init__()
        self.session_maker = session_maker

    @provide(scope=Scope.REQUEST)
    async def session(self) -> AsyncGenerator[AsyncSession]:
        async with self.session_maker() as session:
            try:
                yield session
            except SQLAlchemyError:
                await session.rollback()
            finally:
                await session.commit()


def init_di(config: Settings, session_maker: async_sessionmaker) -> AsyncContainer:
    container = make_async_container(
        DatabaseProvider(session_maker),
        ConfigProvider(config)
    )

    return container
