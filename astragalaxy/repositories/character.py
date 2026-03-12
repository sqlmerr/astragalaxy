from dataclasses import dataclass
from uuid import UUID

from sqlalchemy import select, delete
from sqlalchemy.ext.asyncio import AsyncSession

from astragalaxy.database.models import Character
from astragalaxy.interfaces.character.repo import CharacterRepo


@dataclass(frozen=True)
class CharacterRepository(CharacterRepo):
    session: AsyncSession

    def create_character(self, character: Character) -> None:
        self.session.add(character)

    async def find_one_character(self, id: UUID) -> Character | None:
        stmt = select(Character).where(Character.id == id)
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()

    async def find_one_character_by_code(self, code: str) -> Character | None:
        stmt = select(Character).where(Character.code == code)
        result = await self.session.execute(stmt)
        return result.scalar_one_or_none()

    async def find_all_characters_by_user_id(self, user_id: UUID) -> list[Character]:
        stmt = select(Character).where(Character.user_id == user_id)
        result = await self.session.execute(stmt)

        return list(result.scalars().all())

    async def delete_character(self, id: UUID) -> Character | None:
        stmt = delete(Character).where(Character.id == id)
        await self.session.execute(stmt)

    def update_character(self, character: Character) -> None:
        self.session.add(character)
