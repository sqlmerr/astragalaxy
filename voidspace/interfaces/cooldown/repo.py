from dataclasses import dataclass
from typing import Protocol
from uuid import UUID


@dataclass(frozen=True)
class CreateCooldown:
    character_id: UUID
    seconds: int
    action: str


@dataclass(frozen=True)
class Cooldown:
    set_at: float
    seconds: int
    action: str = "none"


class CooldownRepo(Protocol):
    async def set_cooldown(self, data: CreateCooldown) -> None:
        raise NotImplementedError

    async def get_cooldown(self, character_id: UUID) -> Cooldown:
        raise NotImplementedError
