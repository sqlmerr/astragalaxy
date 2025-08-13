from dataclasses import dataclass

from typing import Protocol
from uuid import UUID

from voidspace.database.models import User


@dataclass(frozen=True)
class UserCreate:
    username: str
    password: str
    token: str


class UserRepo(Protocol):
    async def create_user(self, user: UserCreate) -> UUID:
        raise NotImplementedError

    async def find_one_user(self, id: UUID) -> User | None:
        raise NotImplementedError

    async def find_one_user_by_username(self, username: str) -> User | None:
        raise NotImplementedError

    async def delete_user(self, id: UUID) -> None:
        raise NotImplementedError
