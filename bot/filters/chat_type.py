from typing import Any

from aiogram.enums import ChatType
from aiogram.filters import Filter
from aiogram.types import Message, CallbackQuery


class ChatTypeFilter(Filter):
    def __init__(self, chat_type: ChatType = ChatType.PRIVATE):
        self.chat_type = chat_type

    async def __call__(self, event: Message | CallbackQuery, **kwargs: Any) -> bool | dict[str, Any]:
        if not event.chat or (event.chat and event.chat.type != ChatType.PRIVATE):
            return False

        return True
