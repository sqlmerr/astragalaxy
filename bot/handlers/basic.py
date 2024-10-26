from aiogram.filters import CommandStart
from aiogram.types import Message
from loguru import logger
from aiogram import Router


router = Router()

@router.message(CommandStart())
async def start_cmd(message: Message) -> None:
    await message.answer("Hii")