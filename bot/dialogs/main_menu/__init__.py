from typing import Awaitable, Callable, Any, Coroutine

from aiogram.fsm.state import State
from aiogram.types import CallbackQuery
from aiogram_dialog import Dialog, Window, DialogManager, StartMode
from aiogram_dialog.widgets.kbd import Button
from aiogram_dialog.widgets.text import Const

from .states import MainMenuSG, SpaceshipSG
from . import spaceship
from .. import I18NFormat


def start_dialog(
    state: State,
) -> Callable[[CallbackQuery, Button, DialogManager], Coroutine[Any, Any, None]]:
    async def start_another_dialog(
        callback: CallbackQuery, button: Button, manager: DialogManager
    ):
        await manager.start(state, mode=StartMode.RESET_STACK)

    return start_another_dialog


dialog = Dialog(
    Window(
        I18NFormat("main_menu"),
        Button(
            I18NFormat("spaceship_menu_btn"),
            id="spaceship",
            on_click=start_dialog(SpaceshipSG.choose),
        ),
        state=MainMenuSG.main,
    )
)


def dialogs() -> list[Dialog]:
    return [dialog, spaceship.dialog]
