from dataclasses import dataclass
import random
from uuid import uuid4

from astragalaxy.database.models import Character, Spaceship, Inventory
from astragalaxy.database.models.inventory import InventoryOwnerEnum
from astragalaxy.dto.character import CreateCharacterDTO, CharacterDTO
from astragalaxy.exceptions import AppError
from astragalaxy.exceptions.character import (
    InvalidCharacterCode,
    TooManyCharacters,
    CharacterCodeAlreadyOccupied,
)
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo
from astragalaxy.interfaces.inventory.repo import InventoryRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo
from astragalaxy.interfaces.station.repo import StationRepo
from astragalaxy.interfaces.system.repo import SystemRepo
from astragalaxy.utils import is_valid_username


@dataclass(frozen=True)
class CreateCharacter:
    character_repo: CharacterRepo
    spaceship_repo: SpaceshipRepo
    inventory_repo: InventoryRepo
    system_repo: SystemRepo
    station_repo: StationRepo
    idp: IdentityProvider
    commiter: Commiter

    async def execute(self, data: CreateCharacterDTO) -> CharacterDTO:
        if not is_valid_username(data.code.lower()):
            raise InvalidCharacterCode()

        current_user_id = self.idp.get_current_user_id()
        user_characters = await self.character_repo.find_all_characters_by_user_id(
            current_user_id
        )
        if len(user_characters) >= 5:
            raise TooManyCharacters()

        character = await self.character_repo.find_one_character_by_code(
            data.code.lower()
        )
        if character:
            raise CharacterCodeAlreadyOccupied()

        system = await self.system_repo.find_one_random_system()
        if not system:
            raise AppError

        stations = await self.station_repo.get_stations_by_system(system.id)
        if len(stations) == 0:
            raise AppError
        station = random.choice(stations)

        c = Character(
            id=uuid4(),
            code=data.code.lower(),
            user_id=current_user_id,
            point_id=station.point_id,
        )
        self.character_repo.add(c)

        character_inventory = Inventory(
            id=uuid4(), owner=InventoryOwnerEnum.CHARACTER, owner_id=c.id
        )
        self.inventory_repo.add_inventory(character_inventory)

        sp = Spaceship(
            id=uuid4(),
            name="initial",
            character_id=c.id,
            active=True,
            point_id=station.point_id
        )
        self.spaceship_repo.add(sp)

        spaceship_inventory = Inventory(
            id=uuid4(), owner=InventoryOwnerEnum.SPACESHIP, owner_id=sp.id
        )
        self.inventory_repo.add_inventory(spaceship_inventory)
        character_id = c.id
        await self.commiter.commit()

        character = await self.character_repo.find_one_character(character_id)
        if not character:
            raise AppError()

        return CharacterDTO.from_model(character)
