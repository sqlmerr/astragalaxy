from aiogram.fsm.state import StatesGroup, State


class MainMenuSG(StatesGroup):
    main = State()


class SpaceshipSG(StatesGroup):
    info = State()
