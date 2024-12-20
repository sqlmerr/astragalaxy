from aiogram.fsm.state import StatesGroup, State
from aiogram.types import CallbackQuery
from aiogram_dialog import Dialog, Window, DialogManager
from aiogram_dialog.widgets.kbd import Button, Cancel
from aiogram_dialog.widgets.text import Const
from aiogram_i18n import I18nContext

from dialogs import I18NFormat


class SetLangDialogState(StatesGroup):
    select_language = State()


async def set_lang(
        callback: CallbackQuery, button: Button, dialog_manager: DialogManager
):
    i18n: I18nContext = dialog_manager.middleware_data.get("i18n")
    lang = button.widget_id
    await i18n.set_locale(lang)

    await callback.answer("✅")


dialog = Dialog(
    Window(
        I18NFormat("language_menu"),
        Button(Const("🇺🇸 English"), id="en", on_click=set_lang),
        Button(Const("🇷🇺 Русский"), id="ru", on_click=set_lang),
        Cancel(Const("❌")),
        state=SetLangDialogState.select_language,
    )
)
