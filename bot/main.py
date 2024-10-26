import asyncio

from loguru import logger
from aiogram import Bot, Dispatcher
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode

from api import Api
from api.base import ApiBase
from handlers import basic
from config_reader import config


async def main() -> None:
    api = Api(ApiBase())
    if not (await api.ping()):
        logger.error("Api not responding :(")
        return
    logger.info("Api is working!")

    bot = Bot(config.bot_token, default=DefaultBotProperties(parse_mode=ParseMode.HTML))
    dp = Dispatcher(api=api)
    dp.startup.register(startup)
    dp.include_routers(*[basic.router])

    logger.info("Starting bot...")
    await dp.start_polling(bot)


async def startup() -> None:
    logger.info("Bot successfully started")


if __name__ == "__main__":
    asyncio.run(main())
