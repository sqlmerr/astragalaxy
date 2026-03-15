from typing import Protocol

from astragalaxy.database.models import Point


class PointRepo(Protocol):
    def add(self, point: Point) -> None:
        raise NotImplementedError

    async def find_one_point(self, id: str) -> Point | None:
        raise NotImplementedError

    async def find_all_points_by_system(self, system_id: str) -> list[Point]:
        raise NotImplementedError

    async def delete_point(self, id: str) -> None:
        raise NotImplementedError
