from dataclasses import dataclass
from typing import Any

from redis.asyncio import Redis

from api import Api


KEY_FORMATS = {
    "user_token": "user_token:{user_id}",
    "jwt_token": "jwt_token:{user_id}",
}


@dataclass(frozen=True)
class TokenManager:
    redis: Redis
    api: Api

    async def get_user_token(self, user_id: int) -> str | None:
        value = await self.redis.get(
            KEY_FORMATS["user_token"].format(user_id=str(user_id))
        )
        if value is not None:
            value = f"{user_id}:{value.decode('utf-8')}"
        return value

    async def set_user_token(self, user_id: int, token: str) -> Any:
        token = token.split(":")[1]

        return await self.redis.set(
            KEY_FORMATS["user_token"].format(user_id=str(user_id)), str(token)
        )

    async def get_jwt_token(self, user_id: int) -> str | None:
        value = await self.redis.get(
            KEY_FORMATS["jwt_token"].format(user_id=str(user_id))
        )
        return value

    async def set_jwt_token(self, user_id: int, jwt_token: str) -> Any:
        return await self.redis.set(
            KEY_FORMATS["jwt_token"].format(user_id=str(user_id)), jwt_token
        )
