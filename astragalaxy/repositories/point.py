from dataclasses import dataclass

from sqlalchemy import insert, select, delete
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models import Point
from astragalaxy.interfaces.point.repo import PointRepo


@dataclass(frozen=True)
class PointRepository(PointRepo):
    session: AsyncSession

    def add(self, point: Point) -> None:
        self.session.add(point)

    async def find_one_point(self, id: str) -> Point | None:
        stmt = select(Point).where(Point.id == id)
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()

    async def find_all_points_by_system(self, system_id: str) -> list[Point]:
        stmt = select(Point).where(Point.system_id == system_id)
        result = await self.session.execute(stmt)
        return list(result.scalars().all())

    async def delete_point(self, id: str) -> None:
        stmt = delete(Point).where(Point.id == id)
        await self.session.execute(stmt)
