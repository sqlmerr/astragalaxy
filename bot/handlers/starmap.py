from aiogram import Router
from aiogram.filters import Command
from aiogram.types import Message
from aiogram_dialog import DialogManager, StartMode

from constants import STICKERS
from dialogs.main_menu import StarMapSG

router = Router()


@router.message(Command("starmap"))
async def starmap(message: Message, dialog_manager: DialogManager):
    await message.reply_sticker(STICKERS["starmap_menu"])
    await dialog_manager.start(StarMapSG.select, mode=StartMode.RESET_STACK)
