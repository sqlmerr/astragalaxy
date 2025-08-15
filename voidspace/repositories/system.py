from dataclasses import dataclass

from sqlalchemy import insert, select, delete
from sqlalchemy.ext.asyncio import AsyncSession

from voidspace.database.models import System
from voidspace.interfaces.system.repo import SystemRepo


@dataclass(frozen=True)
class SystemRepository(SystemRepo):
    session: AsyncSession

    async def create_system(self, system: System) -> str:
        stmt = insert(System).values(
            id=system.id, name=system.name, locations=system.locations
        )
        await self.session.execute(stmt)

        return system.id

    async def find_one_system(self, id: str) -> System | None:
        stmt = select(System).where(System.id == id)
        result = await self.session.execute(stmt)

        return result.scalar_one_or_none()

    async def find_all_systems(self, limit: int, offset: int) -> list[System]:
        stmt = select(System).limit(limit).offset(offset)
        result = await self.session.execute(stmt)

        return list(result.scalars().all())

    async def delete_system(self, id: str) -> None:
        stmt = delete(System).where(System.id == id)
        await self.session.execute(stmt)
