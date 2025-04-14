from aiogram.fsm.state import StatesGroup, State

class CreateClient(StatesGroup):
    age = State()
    location = State()
    gender = State()