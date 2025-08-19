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
from redis.asyncio import Redis
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.ext.asyncio import async_sessionmaker, AsyncSession

from voidspace.config import Settings
from voidspace.cooldown_manager import CooldownManager
from voidspace.identity_provider import IdentityProviderImpl
from voidspace.interfaces.character.repo import CharacterRepo
from voidspace.interfaces.cooldown.repo import CooldownRepo
from voidspace.interfaces.identity_provider import IdentityProvider
from voidspace.interfaces.inventory.repo import InventoryRepo
from voidspace.interfaces.item.repo import ItemRepo
from voidspace.interfaces.planet.repo import PlanetRepo
from voidspace.interfaces.resource.repo import ResourceRepo
from voidspace.interfaces.spaceship.repo import SpaceshipRepo
from voidspace.interfaces.system.repo import SystemRepo
from voidspace.interfaces.system_connection.repo import SystemConnectionRepo
from voidspace.interfaces.user.repo import UserRepo
from voidspace.jwt_token_processor import JwtTokenProcessor
from voidspace.password_hasher import PasswordHasher
from voidspace.repositories.character import CharacterRepository
from voidspace.repositories.cooldown import CooldownRepository
from voidspace.repositories.inventory import InventoryRepository
from voidspace.repositories.item import ItemRepository
from voidspace.repositories.planet import PlanetRepository
from voidspace.repositories.resource import ResourceRepository
from voidspace.repositories.spaceship import SpaceshipRepository
from voidspace.repositories.system import SystemRepository
from voidspace.repositories.system_connection import SystemConnectionRepository
from voidspace.repositories.user import UserRepository
from voidspace.use_cases.add_spaceship import AddSpaceship
from voidspace.use_cases.create_character import CreateCharacter
from voidspace.use_cases.create_planet import CreatePlanet
from voidspace.use_cases.create_system import CreateSystem
from voidspace.use_cases.delete_character import DeleteCharacter
from voidspace.use_cases.delete_planet import DeletePlanet
from voidspace.use_cases.delete_system import DeleteSystem
from voidspace.use_cases.enter_spaceship import EnterSpaceship
from voidspace.use_cases.exit_spaceship import ExitSpaceship
from voidspace.use_cases.get_character import GetUserCharacters
from voidspace.use_cases.get_inventory import GetInventoryItems, GetInventoryResources
from voidspace.use_cases.get_planet import GetPlanet, GetSystemPlanets
from voidspace.use_cases.get_spaceship import (
    GetSpaceship,
    GetCharacterSpaceships,
    GetActiveSpaceship,
)
from voidspace.use_cases.get_system import GetSystem, GetSystemsPaginated
from voidspace.use_cases.get_user import GetUserById, GetUserByUsername
from voidspace.use_cases.hyperjump import Hyperjump
from voidspace.use_cases.login import Login
from voidspace.use_cases.navigate_to_planet import NavigateToPlanet
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

    @provide(
        scope=Scope.REQUEST, provides=AnyOf[IdentityProviderImpl, IdentityProvider]
    )
    def get_identity_provider(
        self,
        request: Request,
        jwt_token_processor: JwtTokenProcessor,
        user_repo: UserRepo,
        character_repo: CharacterRepo,
    ):
        return IdentityProviderImpl(
            user_repo=user_repo,
            headers=request.headers,
            jwt_token_processor=jwt_token_processor,
            character_repo=character_repo,
        )

    cooldown_manager = provide(CooldownManager, scope=Scope.APP)


class DatabaseProvider(Provider):
    def __init__(self, session_maker: async_sessionmaker, redis: Redis):
        super().__init__()
        self.session_maker = session_maker
        self.redis = redis

    @provide(scope=Scope.REQUEST)
    async def get_session(self) -> AsyncGenerator[AsyncSession]:
        async with self.session_maker() as session:
            try:
                yield session
            except SQLAlchemyError:
                await session.rollback()
            finally:
                await session.commit()

    @provide(scope=Scope.APP)
    def get_redis(self) -> Redis:
        return self.redis


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

    cooldown_repo = provide(
        CooldownRepository,
        scope=Scope.APP,
        provides=AnyOf[CooldownRepo, CooldownRepository],
    )

    system_connection_repo = provide(
        SystemConnectionRepository,
        scope=Scope.REQUEST,
        provides=AnyOf[SystemConnectionRepo, SystemConnectionRepository],
    )

    inventory_repo = provide(
        InventoryRepository,
        scope=Scope.REQUEST,
        provides=AnyOf[InventoryRepo, InventoryRepository],
    )

    item_repo = provide(
        ItemRepository, scope=Scope.REQUEST, provides=AnyOf[ItemRepo, ItemRepository]
    )

    resource_repo = provide(
        ResourceRepository,
        scope=Scope.REQUEST,
        provides=AnyOf[ResourceRepository, ResourceRepo],
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
        EnterSpaceship,
        ExitSpaceship,
        NavigateToPlanet,
        Hyperjump,
        GetInventoryItems,
        GetInventoryResources,
        scope=Scope.REQUEST,
    )


def init_di(
    config: Settings, session_maker: async_sessionmaker, redis: Redis
) -> AsyncContainer:
    container = make_async_container(
        DatabaseProvider(session_maker, redis),
        CommonProvider(config),
        RepositoryProvider(),
        FastapiProvider(),
        UseCaseProvider(),
    )

    return container
