from aiogram import Router
from aiogram.filters import Command
from aiogram.types import Message
from aiogram_dialog import DialogManager, StartMode

from dialogs.main_menu.states import SpaceshipSG

router = Router()


@router.message(Command("spaceship"))
async def spaceship_cmd(message: Message, dialog_manager: DialogManager) -> None:
    await dialog_manager.start(SpaceshipSG.info, mode=StartMode.RESET_STACK)
