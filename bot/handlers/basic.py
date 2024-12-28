from aiogram.filters import CommandStart, Command
from aiogram.types import Message, ErrorEvent
from aiogram_dialog import DialogManager, StartMode
from aiogram_i18n import I18nContext
from aiogram import Router, F

from api.types.token import TokenPair
from dialogs.main_menu import MainMenuSG
from dialogs.setlang import SetLangDialogState
from filters.admin import AdminFilter
from utils.notifications import notify_admins_error

router = Router()


@router.message(CommandStart())
async def start_cmd(message: Message, dialog_manager: DialogManager) -> None:
    await dialog_manager.start(MainMenuSG.main, mode=StartMode.RESET_STACK)


@router.message(Command("token"))
async def token_cmd(message: Message, i18n: I18nContext, token_pair: TokenPair) -> None:
    await message.reply(i18n.token_menu(token=token_pair.user_token))


@router.message(Command("lang", "setlang", "language", "setlanguage"))
async def set_lang(message: Message, dialog_manager: DialogManager) -> None:
    await dialog_manager.start(
        state=SetLangDialogState.select_language, mode=StartMode.RESET_STACK
    )


@router.message(Command("error"), AdminFilter())
async def raise_error(message: Message) -> None:
    raise ZeroDivisionError()


@router.error(F.update.message.as_("message"))
async def error_handler(
    error: ErrorEvent, message: Message, i18n: I18nContext, **kwargs
) -> None:
    await message.reply(i18n.unexpected_error())

    with i18n.use_locale("ru") as i18n:
        await notify_admins_error(message.bot, error.exception, i18n, message.from_user)
