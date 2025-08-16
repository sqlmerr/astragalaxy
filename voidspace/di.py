from typing import AsyncGenerator

from dishka import (
    Provider,
    make_async_container,
    provide,
    Scope,
    AsyncContainer,
    AnyOf,
    provide_all,
)
from dishka.integrations.fastapi import FastapiProvider
from fastapi import Request
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.ext.asyncio import async_sessionmaker, AsyncSession

from voidspace.config import Settings
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.character.repo import CharacterRepo
from voidspace.interfaces.planet.repo import PlanetRepo
from voidspace.interfaces.spaceship.repo import SpaceshipRepo
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.interfaces.user.repo import UserRepo
from voidspace.jwt_token_processor import JwtTokenProcessor
from voidspace.password_hasher import PasswordHasher
from voidspace.repositories.character import CharacterRepository
from voidspace.repositories.planet import PlanetRepository
from voidspace.repositories.spaceship import SpaceshipRepository
from voidspace.repositories.system import SystemRepository
from voidspace.repositories.user import UserRepository
from voidspace.use_cases.add_spaceship import AddSpaceship
from voidspace.use_cases.create_character import CreateCharacter
from voidspace.use_cases.create_planet import CreatePlanet
from voidspace.use_cases.create_system import CreateSystem
from voidspace.use_cases.delete_character import DeleteCharacter
from voidspace.use_cases.delete_planet import DeletePlanet
from voidspace.use_cases.delete_system import DeleteSystem
from voidspace.use_cases.get_character import GetUserCharacters
from voidspace.use_cases.get_planet import GetPlanet, GetSystemPlanets
from voidspace.use_cases.get_spaceship import (
    GetSpaceship,
    GetCharacterSpaceships,
    GetActiveSpaceship,
)
from voidspace.use_cases.get_system import GetSystem, GetSystemsPaginated
from voidspace.use_cases.get_user import GetUserById, GetUserByUsername
from voidspace.use_cases.login import Login
from voidspace.use_cases.register import Register
from voidspace.use_cases.rename_spaceship import RenameSpaceship
from voidspace.use_cases.set_active_spaceship import SetActiveSpaceship


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
        user_repo: UserRepo,
        character_repo: CharacterRepo,
    ) -> IdentityProvider:
        return IdentityProvider(
            user_repo=user_repo,
            headers=request.headers,
            jwt_token_processor=jwt_token_processor,
            character_repo=character_repo,
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

    system_repo = provide(
        SystemRepository,
        scope=Scope.REQUEST,
        provides=AnyOf[SystemRepo, SystemRepository],
    )

    planet_repo = provide(
        PlanetRepository,
        scope=Scope.REQUEST,
        provides=AnyOf[PlanetRepository, PlanetRepo],
    )

    spaceship_repo = provide(
        SpaceshipRepository,
        scope=Scope.REQUEST,
        provides=AnyOf[SpaceshipRepository, SpaceshipRepo],
    )


class UseCaseProvider(Provider):
    use_cases = provide_all(
        Login,
        Register,
        GetUserById,
        GetUserByUsername,
        CreateCharacter,
        DeleteCharacter,
        GetSystem,
        GetSystemsPaginated,
        CreateSystem,
        DeleteSystem,
        CreatePlanet,
        DeletePlanet,
        GetPlanet,
        GetSystemPlanets,
        GetSpaceship,
        GetCharacterSpaceships,
        GetActiveSpaceship,
        SetActiveSpaceship,
        AddSpaceship,
        RenameSpaceship,
        GetUserCharacters,
        scope=Scope.REQUEST,
    )


def init_di(config: Settings, session_maker: async_sessionmaker) -> AsyncContainer:
    container = make_async_container(
        DatabaseProvider(session_maker),
        CommonProvider(config),
        RepositoryProvider(),
        FastapiProvider(),
        UseCaseProvider(),
    )

    return container
