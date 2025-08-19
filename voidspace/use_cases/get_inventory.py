from dataclasses import dataclass
from uuid import UUID

from voidspace.database.models.inventory import InventoryOwnerEnum
from voidspace.dto.item import ItemDTO
from voidspace.dto.resource import ResourceDTO
from voidspace.exceptions.inventory import InventoryNotFoundError
from voidspace.interfaces.identity_provider import IdentityProvider
from voidspace.interfaces.inventory.repo import InventoryRepo
from voidspace.interfaces.item.repo import ItemRepo
from voidspace.interfaces.resource.repo import ResourceRepo
from voidspace.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class GetInventoryByOwnerRequest:
    owner: InventoryOwnerEnum
    owner_id: UUID


@dataclass(frozen=True)
class GetInventoryByIdRequest:
    inventory_id: UUID


@dataclass(frozen=True)
class GetInventoryItems:
    repo: InventoryRepo
    item_repo: ItemRepo
    spaceship_repo: SpaceshipRepo
    idp: IdentityProvider

    async def execute(
        self,
        data: GetInventoryByOwnerRequest | GetInventoryByIdRequest,
    ) -> list[ItemDTO]:
        current_character = await self.idp.get_current_character()
        if isinstance(data, GetInventoryByOwnerRequest):
            inventory = await self.repo.find_one_inventory_by_owner(
                data.owner, data.owner_id
            )
        else:
            inventory = await self.repo.find_one_inventory(data.inventory_id)

        if not inventory:
            raise InventoryNotFoundError()

        if (
            inventory.owner == InventoryOwnerEnum.CHARACTER
            and inventory.owner_id != current_character.id
        ):
            raise InventoryNotFoundError()

        if inventory.owner == InventoryOwnerEnum.SPACESHIP:
            spaceship = await self.spaceship_repo.find_one_by_id(inventory.owner_id)
            if not spaceship:
                raise InventoryNotFoundError()
            if spaceship.character_id != current_character.id:
                raise InventoryNotFoundError()

        items = await self.item_repo.find_all_items_by_inventory_id(inventory.id)
        return [ItemDTO.from_model(i) for i in items]


@dataclass(frozen=True)
class GetInventoryResources:
    repo: InventoryRepo
    resource_repo: ResourceRepo
    spaceship_repo: SpaceshipRepo
    idp: IdentityProvider

    async def execute(
        self,
        data: GetInventoryByOwnerRequest | GetInventoryByIdRequest,
    ) -> list[ResourceDTO]:
        current_character = await self.idp.get_current_character()
        if isinstance(data, GetInventoryByOwnerRequest):
            inventory = await self.repo.find_one_inventory_by_owner(
                data.owner, data.owner_id
            )
        else:
            inventory = await self.repo.find_one_inventory(data.inventory_id)

        if not inventory:
            raise InventoryNotFoundError()

        if (
            inventory.owner == InventoryOwnerEnum.CHARACTER
            and inventory.owner_id != current_character.id
        ):
            raise InventoryNotFoundError()

        if inventory.owner == InventoryOwnerEnum.SPACESHIP:
            spaceship = await self.spaceship_repo.find_one_by_id(inventory.owner_id)
            if not spaceship:
                raise InventoryNotFoundError()
            if spaceship.character_id != current_character.id:
                raise InventoryNotFoundError()

        resources = await self.resource_repo.find_all_resources_by_inventory_id(
            inventory.id
        )
        return [ResourceDTO.from_model(r) for r in resources]
