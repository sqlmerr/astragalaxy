from typing import Any, reveal_type

from aiogram_dialog.widgets.input import TextInput, ManagedTextInput
from aiogram_i18n import I18nContext
from loguru import logger

from aiogram.types import CallbackQuery, Message
from aiogram_dialog import Dialog, Window, DialogManager, StartMode
from aiogram_dialog.widgets.kbd import Button, SwitchTo
from aiogram_dialog.widgets.text import Const, Format, Case

from api import Api
from api.exceptions import APIError
from api.types.token import TokenPair
from api.types.user import User
from utils.notifications import notify_admins_error
from utils.validators import validate_string
from . import MainMenuSG
from .states import SpaceshipSG
from .. import I18NFormat


async def return_to_main(
    callback: CallbackQuery, button: Button, manager: DialogManager
):
    await manager.start(MainMenuSG.main, mode=StartMode.RESET_STACK)


async def getter(dialog_manager: DialogManager, **kwargs) -> dict:
    api: Api = dialog_manager.middleware_data["api"]
    token_pair: TokenPair = dialog_manager.middleware_data["token_pair"]
    user: User = dialog_manager.middleware_data["user"]

    spaceship = await api.get_my_spaceship(jwt_token=token_pair.jwt_token)

    in_spaceship = dialog_manager.dialog_data["in_spaceship"]

    return {
        "spaceship": spaceship,
        "name": spaceship.name,
        "in_spaceship": in_spaceship,
    }


async def on_start(data: Any, manager: DialogManager):
    user: User = manager.middleware_data["user"]
    manager.dialog_data["in_spaceship"] = user.in_spaceship


async def change_spaceship_status(
    callback: CallbackQuery, button: Button, manager: DialogManager
) -> None:
    api: Api = manager.middleware_data["api"]
    token_pair: TokenPair = manager.middleware_data["token_pair"]
    user: User = manager.middleware_data["user"]

    if user.in_spaceship:
        await api.get_out_of_my_spaceship(token_pair.jwt_token)
        manager.dialog_data["in_spaceship"] = False
    else:
        await api.enter_my_spaceship(token_pair.jwt_token)
        manager.dialog_data["in_spaceship"] = True


async def change_name(message: Message,
            widget: ManagedTextInput,
            dialog_manager: DialogManager,
            data: str):
    i18n: I18nContext = dialog_manager.middleware_data["i18n"]
    api: Api = dialog_manager.middleware_data["api"]
    token_pair: TokenPair = dialog_manager.middleware_data["token_pair"]
    if not validate_string(data):
        await message.reply(i18n.invalid_spaceship_name())
        return
    response = await api.rename_my_spaceship(token_pair.jwt_token, data)
    if response == 2:
        await message.reply(i18n.invalid_spaceship_name())
        return
    if response != 1:
        await notify_admins_error(message.bot, APIError(f"Custom status code = {response}", status_code=response), message.from_user)
        return

    await dialog_manager.switch_to(SpaceshipSG.info)
    await message.reply("✅")


dialog = Dialog(
    Window(
        I18NFormat("spaceship_menu", keys={"name": "name"}),
        Button(
            Case(
                {True: I18NFormat("spaceship_menu_get_out"), False: I18NFormat("spaceship_menu_enter")}, selector="in_spaceship"
            ),
            id="change_spaceship_status",
            on_click=change_spaceship_status,
        ),
        SwitchTo(I18NFormat("spaceship_menu_change_name"), id="change_spaceship_name", state=SpaceshipSG.rename),
        Button(Const("←"), id="to_main", on_click=return_to_main),
        state=SpaceshipSG.info,
        getter=getter,
    ),
    Window(
        I18NFormat("spaceship_menu_enter_name", keys={"name": "name"}),
        TextInput(id="enter_spaceship_name", on_success=change_name),
        SwitchTo(Const("←"), id="to_spaceship", state=SpaceshipSG.info),
        state=SpaceshipSG.rename,
        getter=getter
    ),
    on_start=on_start
)
