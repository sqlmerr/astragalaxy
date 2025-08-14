from typing import AsyncGenerator

from dishka import Provider, make_async_container, provide, Scope, AsyncContainer, AnyOf
from dishka.integrations.fastapi import FastapiProvider
from fastapi import Request
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.ext.asyncio import async_sessionmaker, AsyncSession

from voidspace.config import Settings
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.character.create import CharacterWriter
from voidspace.interfaces.character.delete import CharacterDeleter
from voidspace.interfaces.character.read import CharacterReader
from voidspace.interfaces.character.repo import CharacterRepo
from voidspace.interfaces.user import UserReader, UserWriter
from voidspace.interfaces.user.repo import UserRepo
from voidspace.jwt_token_processor import JwtTokenProcessor
from voidspace.password_hasher import PasswordHasher
from voidspace.repositories.character import CharacterRepository
from voidspace.repositories.user import UserRepository
from voidspace.services.character import CharacterService
from voidspace.services.user import UserService


class CommonProvider(Provider):
    def __init__(self, settings: Settings):
        super().__init__()
        self.settings = settings

    @provide(scope=Scope.APP)
    def get_settings(self) -> Settings:
        return self.settings

    @provide(scope=Scope.APP)
    def get_password_hasher(self) -> PasswordHasher:
        return PasswordHasher()

    jwt_token_processor = provide(JwtTokenProcessor, scope=Scope.APP)

    @provide(scope=Scope.REQUEST)
    def get_identity_provider(
        self,
        request: Request,
        jwt_token_processor: JwtTokenProcessor,
        user_reader: UserReader,
        character_reader: CharacterReader,
    ) -> IdentityProvider:
        return IdentityProvider(
            user_reader=user_reader,
            headers=request.headers,
            jwt_token_processor=jwt_token_processor,
            character_reader=character_reader,
        )


class DatabaseProvider(Provider):
    def __init__(self, session_maker: async_sessionmaker):
        super().__init__()
        self.session_maker = session_maker

    @provide(scope=Scope.REQUEST)
    async def get_session(self) -> AsyncGenerator[AsyncSession]:
        async with self.session_maker() as session:
            try:
                yield session
            except SQLAlchemyError:
                await session.rollback()
            finally:
                await session.commit()


class RepositoryProvider(Provider):
    @provide(scope=Scope.REQUEST, provides=AnyOf[UserRepo, UserRepository])
    async def get_user_repo(self, session: AsyncSession):
        return UserRepository(session)

    character_repo = provide(
        CharacterRepository,
        scope=Scope.REQUEST,
        provides=AnyOf[CharacterRepo, CharacterRepository],
    )


class ServiceProvider(Provider):
    user_service = provide(
        UserService,
        scope=Scope.REQUEST,
        provides=AnyOf[UserReader, UserWriter, UserService],
    )

    character_service = provide(
        CharacterService,
        scope=Scope.REQUEST,
        provides=AnyOf[
            CharacterReader, CharacterWriter, CharacterDeleter, CharacterService
        ],
    )


def init_di(config: Settings, session_maker: async_sessionmaker) -> AsyncContainer:
    container = make_async_container(
        DatabaseProvider(session_maker),
        CommonProvider(config),
        RepositoryProvider(),
        ServiceProvider(),
        FastapiProvider(),
    )

    return container
