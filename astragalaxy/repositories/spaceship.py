from dataclasses import dataclass
from uuid import UUID

from sqlalchemy import select, and_
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models import Spaceship
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class SpaceshipRepository(SpaceshipRepo):
    session: AsyncSession

    def add_spaceship(self, spaceship: Spaceship) -> None:
        self.session.add(spaceship)

    async def find_one_by_id(self, id: UUID) -> Spaceship | None:
        stmt = select(Spaceship).where(Spaceship.id == id)
        result = await self.session.execute(stmt)

        return result.scalar_one_or_none()

    async def find_all_by_character_id(self, character_id: UUID) -> list[Spaceship]:
        stmt = select(Spaceship).where(Spaceship.character_id == character_id)
        result = await self.session.execute(stmt)

        return list(result.scalars().all())

    async def find_one_active_by_character(
        self, character_id: UUID
    ) -> Spaceship | None:
        stmt = select(Spaceship).where(
            and_(Spaceship.character_id == character_id, Spaceship.active == True)
        )
        result = await self.session.execute(stmt)

        return result.scalar_one_or_none()

    def save_spaceship(self, spaceship: Spaceship) -> None:
        self.session.add(spaceship)
