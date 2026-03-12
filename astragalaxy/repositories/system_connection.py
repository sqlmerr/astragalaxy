from dataclasses import dataclass
from uuid import UUID

from sqlalchemy import select, delete, and_
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models.system_connection import SystemConnection
from astragalaxy.interfaces.system_connection.repo import SystemConnectionRepo


@dataclass(frozen=True)
class SystemConnectionRepository(SystemConnectionRepo):
    session: AsyncSession

    def add_connection(self, data: SystemConnection) -> None:
        self.session.add(data)

    async def find_one_connection(self, id: UUID) -> SystemConnection | None:
        stmt = select(SystemConnection).where(SystemConnection.id == id)
        result = await self.session.execute(stmt)

        return result.scalar_one_or_none()

    async def find_all_connections_by_system_id(
        self, system_id: str
    ) -> list[SystemConnection]:
        stmt = select(SystemConnection).where(
            SystemConnection.system_to_id == system_id
        )
        result = await self.session.execute(stmt)

        return list(result.scalars().all())

    async def find_connections_by_system_ids(
        self, system_a_id: str, system_b_id: str
    ) -> list[SystemConnection]:
        stmt = select(SystemConnection).where(
            and_(
                SystemConnection.system_to_id == system_a_id,
                SystemConnection.system_from_id == system_b_id,
            )
            | and_(
                SystemConnection.system_to_id == system_b_id,
                SystemConnection.system_from_id == system_a_id,
            )
        )

        result = await self.session.execute(stmt)
        return list(result.scalars().all())

    async def delete_connection(self, id: UUID) -> None:
        stmt = delete(SystemConnection).where(SystemConnection.id == id)
        await self.session.execute(stmt)
