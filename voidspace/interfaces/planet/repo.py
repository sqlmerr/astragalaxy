from typing import Protocol

from voidspace.database.models import Planet


class PlanetRepo(Protocol):
    async def create_planet(self, planet: Planet) -> str:
        raise NotImplementedError

    async def find_one_planet(self, id: str) -> Planet | None:
        raise NotImplementedError

    async def find_all_planets_by_system(self, system_id: str) -> list[Planet]:
        raise NotImplementedError

    async def delete_planet(self, id: str) -> None:
        raise NotImplementedError
