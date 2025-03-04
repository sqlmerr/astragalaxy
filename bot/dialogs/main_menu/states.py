from aiogram.fsm.state import StatesGroup, State


class MainMenuSG(StatesGroup):
    main = State()


class SpaceshipSG(StatesGroup):
    choose = State()
    info = State()
    rename = State()


class StarMapSG(StatesGroup):
    select = State()
    info = State()
    planet = State()
    flight = State()
