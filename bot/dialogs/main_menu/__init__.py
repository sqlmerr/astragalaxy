from typing import Awaitable, Callable, Any, Coroutine

from aiogram.fsm.state import State
from aiogram.types import CallbackQuery
from aiogram_dialog import Dialog, Window, DialogManager, StartMode
from aiogram_dialog.widgets.kbd import Button

from .states import MainMenuSG, SpaceshipSG, StarMapSG
from . import spaceship, starmap
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
        Button(
            I18NFormat("starmap_menu_btn"),
            id="starmap",
            on_click=start_dialog(StarMapSG.select),
        ),
        state=MainMenuSG.main,
    )
)


def dialogs() -> list[Dialog]:
    return [dialog, spaceship.dialog, starmap.dialog]
