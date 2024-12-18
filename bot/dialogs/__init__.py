from aiogram_dialog import DialogManager
from aiogram_dialog.widgets.common import WhenCondition
from aiogram_dialog.widgets.text import Text
from aiogram_i18n import I18nContext


class I18NFormat(Text):
    def __init__(self, text: str, when: WhenCondition = None, **kwargs):
        super().__init__(when)
        self.text = text
        self.params = kwargs

    async def _render_text(self, data: dict, manager: DialogManager) -> str:
        i18n: I18nContext = manager.middleware_data.get("i18n")
        return i18n.get(self.text, **self.params)
