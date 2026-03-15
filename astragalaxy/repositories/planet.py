from dataclasses import dataclass

from sqlalchemy import insert, select, delete
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models import Planet, Point
from astragalaxy.interfaces.planet.repo import PlanetRepo


@dataclass(frozen=True)
class PlanetRepository(PlanetRepo):
    session: AsyncSession

    def add(self, planet: Planet) -> None:
        self.session.add(planet)

    async def find_one_planet(self, id: str) -> Planet | None:
        stmt = select(Planet).where(Planet.id == id)
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()

    async def get_planets_by_point(self, point_id: str) -> list[Planet]:
        stmt = select(Planet).where(Planet.point_id == point_id)
        result = await self.session.execute(stmt)
        return list(result.scalars().all())

    async def get_planets_by_system(self, system_id: str) -> list[Planet]:
        stmt = (
            select(Planet)
            .join(Planet.point)
            .where(Point.system_id == system_id)
        )

        result = await self.session.scalars(stmt)
        return list(result.all())

    async def delete_planet(self, id: str) -> None:
        stmt = delete(Planet).where(Planet.id == id)
        await self.session.execute(stmt)
