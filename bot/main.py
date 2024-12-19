import asyncio

from aiogram.fsm.storage.base import DefaultKeyBuilder
from aiogram.types import BotCommand, BotCommandScopeDefault
from aiogram_dialog import setup_dialogs
from aiogram_i18n import I18nMiddleware
from aiogram_i18n.cores import FluentRuntimeCore
from loguru import logger
from aiogram import Bot, Dispatcher
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode
from aiogram.fsm.storage.redis import RedisStorage
from redis.asyncio import Redis

from dialogs import setlang
from api import Api
from api.base import ApiBase
from handlers import basic
from config_reader import config


async def main() -> None:
    api = Api(ApiBase())
    if not (await api.ping()):
        logger.error("Api isn't responding :(\nCan't start application")
        return
    logger.info("Api is working!")

    redis = Redis.from_url(config.redis_url)
    storage = RedisStorage(redis, key_builder=DefaultKeyBuilder(with_destiny=True, with_bot_id=True))

    bot = Bot(config.bot_token, default=DefaultBotProperties(parse_mode=ParseMode.HTML))
    dp = Dispatcher(api=api, storage=storage, redis=redis)
    dp.startup.register(startup)
    dp.include_routers(*[basic.router, setlang.dialog])
    setup_dialogs(router=dp)

    i18n_middleware = I18nMiddleware(
        core=FluentRuntimeCore(
            path="locales/{locale}",
            default_locale="ru"
        ),
        default_locale="ru"
    )
    i18n_middleware.setup(dp)

    me = await bot.me()
    logger.info(f"Starting bot @{me.username}...")
    await bot.delete_webhook(drop_pending_updates=True)
    await set_commands(bot)
    await dp.start_polling(bot)


async def startup() -> None:
    logger.info("Bot successfully started")


async def set_commands(bot: Bot) -> None:
    commands = {
        "start": "Restart bot",
        "lang": "Changle language"
    }

    cmds = [BotCommand(command=cmd, description=desc) for cmd, desc in commands.items()]

    await bot.set_my_commands(cmds, scope=BotCommandScopeDefault())


if __name__ == "__main__":
    asyncio.run(main())
