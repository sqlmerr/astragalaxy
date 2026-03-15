from typing import Protocol

from astragalaxy.database.models import Station


class StationRepo(Protocol):
    def add(self, station: Station) -> None:
        raise NotImplementedError

    async def find_one_station(self, id: str) -> Station | None:
        raise NotImplementedError

    async def get_stations_by_point(self, point_id: str) -> list[Station]:
        raise NotImplementedError

    async def get_stations_by_system(self, system_id: str) -> list[Station]:
        raise NotImplementedError

    async def delete_station(self, id: str) -> None:
        raise NotImplementedError
