from aiogram.fsm.state import StatesGroup, State

class CreateAdvertiser(StatesGroup):
    name = State()

class CreateCampaign(StatesGroup):
    name = State()
    ad_title = State()
    ad_text = State()
    impressions_limit = State()
    clicks_limit = State()
    cost_per_impression = State()
    cost_per_click = State()
    start_date = State()
    end_date = State()
    gender = State()
    age_from = State()
    age_to = State()
    location = State()