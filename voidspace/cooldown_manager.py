import time
from dataclasses import dataclass
from uuid import UUID

from voidspace.dto.cooldown import SetCooldownDTO, CooldownDTO
from voidspace.interfaces.cooldown.repo import CooldownRepo, CreateCooldown


def _get_remaining_seconds(set_at: float, seconds: float) -> int:
    now = time.time()  # 1050

    remaining = set_at + seconds - int(now)
    return int(remaining if remaining > 0 else 0)


@dataclass(frozen=True)
class CooldownManager:
    repo: CooldownRepo

    async def set(self, data: SetCooldownDTO) -> CooldownDTO:
        await self.repo.set_cooldown(
            CreateCooldown(
                character_id=data.character_id, seconds=data.seconds, action=data.action
            )
        )

        return await self.get(data.character_id)

    async def get(self, character_id: UUID) -> CooldownDTO:
        cooldown = await self.repo.get_cooldown(character_id)
        return CooldownDTO(
            set_at=cooldown.set_at,  # 1000
            action=cooldown.action,
            seconds=cooldown.seconds,  # 30
            remaining_seconds=_get_remaining_seconds(
                cooldown.set_at, seconds=cooldown.seconds
            ),
        )
