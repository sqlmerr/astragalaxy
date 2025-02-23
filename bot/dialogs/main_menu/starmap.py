from aiogram.types import CallbackQuery
from aiogram_dialog import Dialog, Window, DialogManager, StartMode
from aiogram_dialog.widgets.kbd import (
    Button,
    StubScroll,
    Row,
    NextPage,
    LastPage,
    CurrentPage,
    PrevPage,
    FirstPage,
    ListGroup, SwitchTo,
)
from aiogram_dialog.widgets.text import Const, Format

from api import Api
from api.types.system import System
from api.types.token import TokenPair
from dialogs import I18NFormat
from dialogs.main_menu import SpaceshipSG
from dialogs.main_menu.states import StarMapSG, MainMenuSG


ID_SCROLL = "scroll"
ITEMS_PER_PAGE = 3


async def return_to_main(
    callback: CallbackQuery, button: Button, manager: DialogManager
):
    await manager.start(MainMenuSG.main, mode=StartMode.RESET_STACK)


async def open_system_info(
    callback: CallbackQuery, button: Button, manager: DialogManager
) -> None:
    system_id = manager.item_id
    api: Api = manager.middleware_data["api"]
    token_pair: TokenPair = manager.middleware_data["token_pair"]
    system = await api.get_system_by_id(token_pair.jwt_token, system_id)

    manager.dialog_data["system"] = system.model_dump(mode="json")
    await manager.switch_to(StarMapSG.info)


async def paging_getter(dialog_manager: DialogManager, **_kwargs):
    api: Api = dialog_manager.middleware_data["api"]
    token_pair: TokenPair = dialog_manager.middleware_data["token_pair"]

    systems = await api.get_all_systems(token_pair.jwt_token)

    current_page = await dialog_manager.find(ID_SCROLL).get_page()
    content = (
        systems[current_page * ITEMS_PER_PAGE :]
        if len(systems) < ITEMS_PER_PAGE
        else systems[
            current_page * ITEMS_PER_PAGE : current_page * ITEMS_PER_PAGE
            + ITEMS_PER_PAGE
        ]
    )
    return {
        "pages": (len(systems) // ITEMS_PER_PAGE) + 1,
        "current_page": current_page + 1,
        "content": content,
    }


async def info_getter(dialog_manager: DialogManager, **_kwargs):
    api: Api = dialog_manager.middleware_data["api"]
    token_pair: TokenPair = dialog_manager.middleware_data["token_pair"]
    system: System = System.model_validate(dialog_manager.dialog_data["system"])

    planets = await api.get_system_planets(token_pair.jwt_token, system.id)

    return {
        "system": system,
        "name": system.name,
        "id": str(system.id),
        "planets": planets,
    }


dialog = Dialog(
    Window(
        I18NFormat("starmap_menu"),
        Format("<b>{current_page}/{pages}</b>"),
        ListGroup(
            Button(Format("{item.name}"), id="system", on_click=open_system_info),
            item_id_getter=lambda i: i.id,
            items="content",
            id="listgroup",
        ),
        StubScroll(id=ID_SCROLL, pages="pages"),
        Row(
            FirstPage(
                scroll=ID_SCROLL,
                text=Format("⏮️ "),
            ),
            PrevPage(
                scroll=ID_SCROLL,
                text=Format("◀️"),
            ),
            CurrentPage(
                scroll=ID_SCROLL,
                text=Format("{current_page1}"),
            ),
            NextPage(
                scroll=ID_SCROLL,
                text=Format("▶️"),
            ),
            LastPage(
                scroll=ID_SCROLL,
                text=Format(" ⏭️"),
            ),
        ),
        Button(Const("←"), id="to_main", on_click=return_to_main),
        getter=paging_getter,
        state=StarMapSG.select,
    ),
    Window(
        I18NFormat("starmap_system_info", keys={"name": "name", "id": "id"}),
        ListGroup(
            Button(Format("{item.name}"), id="planet"),
            id="planet_list",
            items="planets",
            item_id_getter=lambda i: i.id,
        ),
        SwitchTo(Const("←"), id="to_select", state=StarMapSG.select),
        state=StarMapSG.info,
        getter=info_getter,
    ),
)
