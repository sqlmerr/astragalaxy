import asyncio

from aiogram import Bot, Dispatcher
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode
from aiogram.filters import ExceptionTypeFilter
from aiogram.fsm.storage.base import DefaultKeyBuilder
from aiogram.fsm.storage.redis import RedisStorage
from aiogram.types import BotCommand, BotCommandScopeDefault, ErrorEvent
from aiogram_dialog import setup_dialogs
from aiogram_dialog.api.exceptions import UnknownIntent, OutdatedIntent
from aiogram_i18n import I18nMiddleware
from aiogram_i18n.cores import FluentRuntimeCore
from loguru import logger
from redis.asyncio import Redis

from api import Api
from api.base import ApiBase
from config_reader import config
from dialogs import setlang, main_menu
from handlers import basic, spaceship, starmap
from middlewares import UserMiddleware
from utils.token_manager import TokenManager


async def main() -> None:
    api = Api(ApiBase())
    if not (await api.ping()):
        logger.error("Api isn't responding :(\nCan't start application")
        return
    logger.info("Api is working!")

    redis = Redis.from_url(config.redis_url)
    storage = RedisStorage(
        redis, key_builder=DefaultKeyBuilder(with_destiny=True, with_bot_id=True)
    )

    token_manager = TokenManager(redis=redis, api=api)

    bot = Bot(config.bot_token, default=DefaultBotProperties(parse_mode=ParseMode.HTML))
    dp = Dispatcher(api=api, storage=storage, redis=redis, token_manager=token_manager)
    dp.startup.register(startup)
    dp.error.register(
        dialog_errors_handler, ExceptionTypeFilter(UnknownIntent, OutdatedIntent)
    )
    dp.include_routers(basic.router, spaceship.router, starmap.router)
    dp.include_routers(setlang.dialog, *main_menu.dialogs())
    setup_dialogs(router=dp)
    dp.message.middleware(UserMiddleware())
    dp.callback_query.middleware(UserMiddleware())

    i18n_middleware = I18nMiddleware(
        core=FluentRuntimeCore(path="locales/{locale}", default_locale="ru"),
        default_locale="ru",
    )
    i18n_middleware.setup(dp)

    me = await bot.me()
    logger.info(f"Starting bot @{me.username}...")
    await bot.delete_webhook(drop_pending_updates=True)
    await set_commands(bot)
    await dp.start_polling(bot)


async def startup() -> None:
    logger.info("Bot successfully started")


async def dialog_errors_handler(error: ErrorEvent):
    update = error.update
    if hasattr(update, "callback_query"):
        await update.callback_query.answer()
        await update.callback_query.message.edit_reply_markup(reply_markup=None)


async def set_commands(bot: Bot) -> None:
    commands = {
        "start": "main menu",
        "spaceship": "your spaceship",
        "starmap": "star map",
        "lang": "change language",
    }

    cmds = [BotCommand(command=cmd, description=desc) for cmd, desc in commands.items()]

    await bot.set_my_commands(cmds, scope=BotCommandScopeDefault())


if __name__ == "__main__":
    asyncio.run(main())
