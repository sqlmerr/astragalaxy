from aiogram.types import CallbackQuery, InlineKeyboardButton, CopyTextButton
from aiogram_dialog import DialogProtocol, DialogManager
from aiogram_dialog.api.internal import RawKeyboard
from aiogram_dialog.widgets.common import WhenCondition
from aiogram_dialog.widgets.kbd import Keyboard
from aiogram_dialog.widgets.kbd.button import OnClick
from aiogram_dialog.widgets.text import Text
from aiogram_dialog.widgets.widget_event import WidgetEventProcessor


class CopyButton(Keyboard):
    def __init__(
        self,
        text: Text,
        id: str,
        copy_text: str | None = None,
        copy_text_key: str | None = None,
        when: WhenCondition = None,
    ):
        super().__init__(id=id, when=when)
        self.text = text
        self.copy_text = copy_text
        self.copy_text_key = copy_text_key

    async def _render_keyboard(
        self,
        data: dict,
        manager: DialogManager,
    ) -> RawKeyboard:
        copy_text = self.copy_text or data.get(self.copy_text_key, "")

        return [
            [
                InlineKeyboardButton(
                    text=await self.text.render_text(data, manager),
                    copy_text=CopyTextButton(text=copy_text),
                ),
            ],
        ]
