from dataclasses import dataclass
from uuid import uuid4

from astragalaxy.database.models import Spaceship, Inventory
from astragalaxy.database.models.inventory import InventoryOwnerEnum
from astragalaxy.dto.spaceship import CreateSpaceshipDTO, SpaceshipDTO
from astragalaxy.exceptions import AppError
from astragalaxy.exceptions.spaceship import TooManySpaceshipsError
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.inventory.repo import InventoryRepo
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class AddSpaceship:
    repo: SpaceshipRepo
    inventory_repo: InventoryRepo
    idp: IdentityProvider

    async def execute(self, data: CreateSpaceshipDTO) -> SpaceshipDTO:
        character = await self.idp.get_current_character()

        character_spaceships = await self.repo.find_all_by_character_id(character.id)
        if len(character_spaceships) >= 3:
            raise TooManySpaceshipsError()
        sp = Spaceship(
            id=uuid4(),
            name=data.name,
            location="space_station",
            character_id=character.id,
            active=False,
            system_id=data.system_id,
            planet_id=None,
        )
        self.repo.add_spaceship(sp)
        spaceship_inventory = Inventory(
            id=uuid4(), owner=InventoryOwnerEnum.SPACESHIP, owner_id=sp.id
        )
        self.inventory_repo.add_inventory(spaceship_inventory)

        spaceship = await self.repo.find_one_by_id(sp.id)
        if not spaceship:
            raise AppError()
        return SpaceshipDTO.from_model(spaceship)
