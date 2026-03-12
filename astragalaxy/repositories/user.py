from dataclasses import dataclass
from uuid import UUID

from sqlalchemy import select, insert, delete
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models import User
from astragalaxy.interfaces.user.repo import UserRepo, UserCreate


@dataclass(frozen=True)
class UserRepository(UserRepo):
    session: AsyncSession

    async def create_user(self, user: UserCreate) -> UUID:
        stmt = (
            insert(User)
            .values(username=user.username, password=user.password, token=user.token)
            .returning(User.id)
        )
        result = await self.session.execute(stmt)

        return result.scalar_one()

    async def find_one_user(self, id: UUID) -> User | None:
        stmt = select(User).where(User.id == id)
        result = await self.session.execute(stmt)

        return result.scalar_one_or_none()

    async def find_one_user_by_username(self, username: str) -> User | None:
        stmt = select(User).where(User.username == username)
        result = await self.session.execute(stmt)

        return result.scalar_one_or_none()

    async def delete_user(self, id: UUID) -> None:
        stmt = delete(User).where(User.id == id)
        await self.session.execute(stmt)
