from typing import Protocol

from astragalaxy.database.models import Planet


class PlanetRepo(Protocol):
    def add(self, planet: Planet) -> None:
        raise NotImplementedError

    async def find_one_planet(self, id: str) -> Planet | None:
        raise NotImplementedError

    async def get_planets_by_point(self, point_id: str) -> list[Planet]:
        raise NotImplementedError

    async def get_planets_by_system(self, system_id: str) -> list[Planet]:
        raise NotImplementedError

    async def delete_planet(self, id: str) -> None:
        raise NotImplementedError
