from typing import Any, reveal_type

from aiogram_dialog.widgets.input import TextInput, ManagedTextInput
from aiogram_i18n import I18nContext
from loguru import logger

from aiogram.types import CallbackQuery, Message
from aiogram_dialog import Dialog, Window, DialogManager, StartMode
from aiogram_dialog.widgets.kbd import Button, SwitchTo, ListGroup
from aiogram_dialog.widgets.text import Const, Format, Case

from api import Api
from api.exceptions import APIError
from api.types.spaceship import Spaceship
from api.types.token import TokenPair
from api.types.user import User
from utils.notifications import notify_admins_error
from utils.validators import validate_string
from . import MainMenuSG
from .states import SpaceshipSG
from .. import I18NFormat
from ..widgets.copy_button import CopyButton


async def return_to_main(
    callback: CallbackQuery, button: Button, manager: DialogManager
):
    await manager.start(MainMenuSG.main, mode=StartMode.RESET_STACK)


async def getter(dialog_manager: DialogManager, **kwargs) -> dict:
    api: Api = dialog_manager.middleware_data["api"]
    token_pair: TokenPair = dialog_manager.middleware_data["token_pair"]
    spaceship: Spaceship = Spaceship.model_validate(dialog_manager.dialog_data["spaceship"])

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
    spaceship: Spaceship = Spaceship.model_validate(manager.dialog_data["spaceship"])

    if user.in_spaceship:
        await api.exit_my_spaceship(token_pair.jwt_token, spaceship.id)
        manager.dialog_data["in_spaceship"] = False
    else:
        await api.enter_my_spaceship(token_pair.jwt_token, spaceship.id)
        manager.dialog_data["in_spaceship"] = True


async def change_name(message: Message,
            widget: ManagedTextInput,
            dialog_manager: DialogManager,
            data: str):
    i18n: I18nContext = dialog_manager.middleware_data["i18n"]
    api: Api = dialog_manager.middleware_data["api"]
    token_pair: TokenPair = dialog_manager.middleware_data["token_pair"]
    spaceship: Spaceship = Spaceship.model_validate(dialog_manager.dialog_data["spaceship"])
    response = await api.rename_my_spaceship(token_pair.jwt_token, spaceship.id, data)
    if response == 2:
        await message.reply(i18n.invalid_spaceship_name())
        return
    if response != 1:
        await notify_admins_error(message.bot, APIError(f"Custom status code = {response}", status_code=response), message.from_user)
        return

    spaceships = await api.get_my_spaceships(token_pair.jwt_token)
    spaceship_ = list(filter(lambda s: s.id == spaceship.id, spaceships))
    if not spaceship_:
        return
    spaceship = spaceship_[0]
    dialog_manager.dialog_data["spaceship"] = spaceship.model_dump(mode="json")

    await dialog_manager.switch_to(SpaceshipSG.info)
    await message.reply("✅")


async def spaceships(dialog_manager: DialogManager, **kwargs) -> dict:
    user: User = dialog_manager.middleware_data["user"]

    spaceships = [{"key": i, "spaceship": s} for i, s in enumerate(user.spaceships)]
    return {
        "spaceships": spaceships,
    }


async def select_spaceship(callback: CallbackQuery, button: Button, manager: DialogManager):
    user: User = manager.middleware_data["user"]
    spaceship = list(filter(lambda s: str(s.id) == manager.item_id, user.spaceships))
    if not spaceship:
        await callback.answer()
        return
    spaceship = spaceship[0]
    manager.dialog_data["spaceship"] = spaceship.model_dump(mode="json")
    await manager.switch_to(SpaceshipSG.info)


dialog = Dialog(
    Window(
        I18NFormat("spaceship_menu_choose_spaceship"),
        ListGroup(
            Button(
                Format("{item[key]}. {item[spaceship].name}"),
                on_click=select_spaceship,
                id="s",
            ),
            id="choose_spaceship",
            item_id_getter=lambda s: s["spaceship"].id,
            items="spaceships"
        ),
        Button(Const("←"), id="to_main", on_click=return_to_main),
        state=SpaceshipSG.choose,
        getter=spaceships
    ),
    Window(
        I18NFormat("spaceship_menu", keys={"name": "name"}),
        Button(
            Case(
                {True: I18NFormat("spaceship_menu_exit"), False: I18NFormat("spaceship_menu_enter")}, selector="in_spaceship"
            ),
            id="change_spaceship_status",
            on_click=change_spaceship_status,
        ),
        SwitchTo(I18NFormat("spaceship_menu_change_name"), id="change_spaceship_name", state=SpaceshipSG.rename),
        SwitchTo(Const("←"), id="to_spaceship", state=SpaceshipSG.choose),
        state=SpaceshipSG.info,
        getter=getter,
    ),
    Window(
        I18NFormat("spaceship_menu_enter_name", keys={"name": "name"}),
        TextInput(id="enter_spaceship_name", on_success=change_name),
        CopyButton(I18NFormat("spaceship_menu_copy_name", keys={"name": "name"}), id="copy_spaceship_name", copy_text_key="name"),
        SwitchTo(Const("←"), id="to_spaceship", state=SpaceshipSG.info),
        state=SpaceshipSG.rename,
        getter=getter
    ),
    on_start=on_start
)
