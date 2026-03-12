from dataclasses import dataclass
from uuid import uuid4

from astragalaxy.database.models import Character, Spaceship, Inventory
from astragalaxy.database.models.inventory import InventoryOwnerEnum
from astragalaxy.dto.character import CreateCharacterDTO, CharacterDTO
from astragalaxy.exceptions import AppError
from astragalaxy.exceptions.character import (
    TooManyCharacters,
    CharacterCodeAlreadyOccupied,
)
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo
from astragalaxy.interfaces.inventory.repo import InventoryRepo
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo
from astragalaxy.interfaces.system.repo import SystemRepo


@dataclass(frozen=True)
class CreateCharacter:
    character_repo: CharacterRepo
    spaceship_repo: SpaceshipRepo
    inventory_repo: InventoryRepo
    system_repo: SystemRepo
    idp: IdentityProvider

    async def execute(self, data: CreateCharacterDTO) -> CharacterDTO:
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

        c = Character(
            id=uuid4(),
            code=data.code.lower(),
            location="space_station",
            user_id=current_user_id,
            system_id=system.id,
        )
        self.character_repo.create_character(c)

        character_inventory = Inventory(
            id=uuid4(), owner=InventoryOwnerEnum.CHARACTER, owner_id=c.id
        )
        self.inventory_repo.add_inventory(character_inventory)

        sp = Spaceship(
            id=uuid4(),
            name="initial",
            location="space_station",
            character_id=c.id,
            active=True,
            system_id=system.id,
            planet_id=None,
        )
        self.spaceship_repo.add_spaceship(sp)

        spaceship_inventory = Inventory(
            id=uuid4(), owner=InventoryOwnerEnum.SPACESHIP, owner_id=sp.id
        )
        self.inventory_repo.add_inventory(spaceship_inventory)

        character = await self.character_repo.find_one_character(c.id)
        if not character:
            raise AppError()

        return CharacterDTO.from_model(character)
