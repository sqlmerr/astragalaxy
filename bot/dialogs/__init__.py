from aiogram_dialog import DialogManager
from aiogram_dialog.widgets.common import WhenCondition
from aiogram_dialog.widgets.text import Text
from aiogram_i18n import I18nContext


class I18NFormat(Text):
    def __init__(
        self,
        text: str,
        keys: dict[str, str] | None = None,
        when: WhenCondition = None,
        **kwargs,
    ):
        super().__init__(when)
        self.text = text
        self.keys = keys
        self.params = kwargs

    async def _render_text(self, data: dict, manager: DialogManager) -> str:
        i18n: I18nContext = manager.middleware_data.get("i18n")
        if not self.keys:
            text = i18n.get(self.text, **self.params)
        else:
            params = {}
            for param, key in self.keys.items():
                params[param] = data.get(key, "undefined")
            text = i18n.get(self.text, **params, **self.params)

        return text
