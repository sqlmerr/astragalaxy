from typing import Callable, Awaitable, Any

from aiogram import BaseMiddleware
from aiogram.types import Message


class UserMiddleware(BaseMiddleware):
    async def __call__(
        self,
        handler: Callable[[Message, dict[str, Any]], Awaitable[Any]],
        event: Message,
        data: dict[str, Any],
    ) -> Any:
        user_id = event.from_user.id

        return await handler(event, data)
