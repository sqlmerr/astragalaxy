from dataclasses import dataclass
from uuid import UUID


@dataclass(frozen=True)
class CooldownDTO:
    set_at: float
    seconds: int
    remaining_seconds: int
    action: str


@dataclass(frozen=True)
class SetCooldownDTO:
    character_id: UUID
    seconds: int
    action: str
