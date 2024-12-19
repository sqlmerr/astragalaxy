from aiogram.filters import CommandStart, Command
from aiogram.types import Message, ErrorEvent, LinkPreviewOptions
from aiogram_dialog import DialogManager, setup_dialogs, StartMode
from aiogram_i18n import I18nContext
from loguru import logger
from aiogram import Router, F

from api import Api
from config_reader import config
from dialogs.setlang import dialog, SetLangDialogState

router = Router()


@router.message(CommandStart())
async def start_cmd(message: Message, api: Api) -> None:
    await api.register_user(message.from_user.id, message.from_user.username)

    await message.answer("Hii")


@router.message(Command("lang", "setlang", "language", "setlanguage"))
async def set_lang(message: Message, dialog_manager: DialogManager) -> None:
    await dialog_manager.start(state=SetLangDialogState.select_language, mode=StartMode.RESET_STACK)


@router.error(F.update.message.as_("message"))
async def error_handler(
    error: ErrorEvent, message: Message, i18n: I18nContext, **kwargs
) -> None:
    await message.reply(i18n.unexpected_error())

    for admin in config.admins:
        with i18n.use_locale(
            await i18n.manager.locale_getter(event=message, **kwargs)
        ) as i18n:
            await message.bot.send_message(
                admin,
                text=i18n.error_admin_notification(
                    user=f"<a href='tg://user?id={message.from_user.id}'>{message.from_user.first_name}</a>",
                    error=str(error.exception),
                ),
                link_preview_options=LinkPreviewOptions(is_disabled=True),
            )
