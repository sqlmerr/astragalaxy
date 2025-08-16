from dataclasses import dataclass

from sqlalchemy import insert, select, delete
from sqlalchemy.ext.asyncio import AsyncSession

from voidspace.database.models import Planet
from voidspace.interfaces.planet.repo import PlanetRepo


@dataclass(frozen=True)
class PlanetRepository(PlanetRepo):
    session: AsyncSession

    async def create_planet(self, planet: Planet) -> str:
        stmt = insert(Planet).values(
            id=planet.id,
            name=planet.name,
            system_id=planet.system_id,
            threat=planet.threat,
        )
        await self.session.execute(stmt)
        return planet.id

    async def find_one_planet(self, id: str) -> Planet | None:
        stmt = select(Planet).where(Planet.id == id)
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()

    async def find_all_planets_by_system(self, system_id: str) -> list[Planet]:
        stmt = select(Planet).where(Planet.system_id == system_id)
        result = await self.session.execute(stmt)
        return list(result.scalars().all())

    async def delete_planet(self, id: str) -> None:
        stmt = delete(Planet).where(Planet.id == id)
        await self.session.execute(stmt)
