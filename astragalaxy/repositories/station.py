from dataclasses import dataclass

from sqlalchemy import insert, select, delete
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models import Point, Station
from astragalaxy.interfaces.station.repo import StationRepo


@dataclass(frozen=True)
class StationRepository(StationRepo):
    session: AsyncSession

    def add(self, station: Station) -> None:
        self.session.add(station)

    async def find_one_station(self, id: str) -> Station | None:
        stmt = select(Station).where(Station.id == id)
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()

    async def get_stations_by_point(self, point_id: str) -> list[Station]:
        stmt = select(Station).where(Station.point_id == point_id)
        result = await self.session.execute(stmt)
        return list(result.scalars().all())

    async def get_stations_by_system(self, system_id: str) -> list[Station]:
        stmt = (
            select(Station)
            .join(Station.point)
            .where(Point.system_id == system_id)
        )

        result = await self.session.scalars(stmt)
        return list(result.all())

    async def delete_station(self, id: str) -> None:
        stmt = delete(Station).where(Station.id == id)
        await self.session.execute(stmt)
