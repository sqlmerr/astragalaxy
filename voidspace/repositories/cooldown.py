import json
import time
from dataclasses import dataclass
from uuid import UUID

from redis.asyncio import Redis

from voidspace.interfaces.cooldown.repo import CooldownRepo, CreateCooldown, Cooldown


@dataclass(frozen=True)
class CooldownRepository(CooldownRepo):
    redis: Redis

    async def set_cooldown(self, data: CreateCooldown) -> None:
        cooldown = {
            "seconds": data.seconds,
            "action": data.action,
            "set_at": time.time(),
        }

        await self.redis.set(str(data.character_id), json.dumps(cooldown))

    async def get_cooldown(self, character_id: UUID) -> Cooldown:
        raw = await self.redis.get(str(character_id))
        if raw is None:
            return Cooldown(seconds=0, set_at=0)
        cooldown = json.loads(raw)
        return Cooldown(
            seconds=cooldown.get("seconds"),
            action=cooldown.get("action"),
            set_at=cooldown.get("set_at"),
        )
