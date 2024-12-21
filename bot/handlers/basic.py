from aiogram.filters import CommandStart, Command
from aiogram.types import Message, ErrorEvent, LinkPreviewOptions
from aiogram_dialog import DialogManager, setup_dialogs, StartMode, ShowMode
from aiogram_i18n import I18nContext
from loguru import logger
from aiogram import Router, F

from api import Api
from api.exceptions import APIError
from api.types.token import TokenPair
from api.types.user import User
from config_reader import config
from dialogs.setlang import dialog, SetLangDialogState
from filters.admin import AdminFilter
from utils.notifications import notify_admins_error

router = Router()


@router.message(CommandStart())
async def start_cmd(message: Message) -> None:
    await message.answer("Hii")


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
    error: ErrorEvent, message: Message, i18n: I18nContext,**kwargs
) -> None:
    await message.reply(i18n.unexpected_error())

    with i18n.use_locale(
        "ru"
    ) as i18n:
        await notify_admins_error(message.bot, error.exception, i18n, message.from_user)
