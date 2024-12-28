from typing import Any

from aiogram.filters import Filter
from aiogram.types import Message, CallbackQuery

from config_reader import config


class AdminFilter(Filter):
    async def __call__(self, event: Message | CallbackQuery) -> bool | dict[str, Any]:
        if event.from_user.id not in config.admins:
            return False
        return True
