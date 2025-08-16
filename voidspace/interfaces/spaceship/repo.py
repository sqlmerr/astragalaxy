from typing import Protocol
from uuid import UUID

from voidspace.database.models import Spaceship


class SpaceshipRepo(Protocol):
    def add_spaceship(self, spaceship: Spaceship) -> None:
        raise NotImplementedError

    async def find_one_by_id(self, id: UUID) -> Spaceship | None:
        raise NotImplementedError

    async def find_all_by_character_id(self, character_id: UUID) -> list[Spaceship]:
        raise NotImplementedError

    async def find_one_active_by_character(
        self, character_id: UUID
    ) -> Spaceship | None:
        raise NotImplementedError

    def save_spaceship(self, spaceship: Spaceship) -> None:
        raise NotImplementedError
