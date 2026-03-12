from dataclasses import dataclass
from uuid import UUID

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models import Resource
from astragalaxy.interfaces.resource.repo import ResourceRepo


@dataclass(frozen=True)
class ResourceRepository(ResourceRepo):
    session: AsyncSession

    def add_resource(self, resource: Resource) -> None:
        self.session.add(resource)

    async def find_one_resource(self, id: UUID) -> Resource | None:
        stmt = select(Resource).where(Resource.id == id)
        result = await self.session.execute(stmt)

        return result.scalar_one_or_none()

    async def find_one_resource_by_code_and_inventory(
        self, code: str, inventory_id: UUID
    ) -> Resource | None:
        stmt = select(Resource).where(
            (Resource.inventory_id == inventory_id) & (Resource.code == code)
        )
        result = await self.session.execute(stmt)

        return result.scalar_one_or_none()

    async def find_all_resources_by_inventory_id(
        self, inventory_id: UUID
    ) -> list[Resource]:
        stmt = select(Resource).where(Resource.inventory_id == inventory_id)
        result = await self.session.execute(stmt)

        return list(result.scalars().all())
