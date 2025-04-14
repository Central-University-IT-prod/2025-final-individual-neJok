from config import BACK_BUTTON, SKIP_BUTTON
from aiogram.types import ReplyKeyboardMarkup, KeyboardButton, InlineKeyboardButton, InlineKeyboardMarkup

kb_choose_profile = InlineKeyboardMarkup(
    inline_keyboard=[
        [
            InlineKeyboardButton(text="Клиента", callback_data="client_profile"),
            InlineKeyboardButton(text="Рекламодателя", callback_data="advertiser_profile")
        ],
    ],
)

kb_back = ReplyKeyboardMarkup(
    keyboard=[
        [
            KeyboardButton(text=BACK_BUTTON)
        ],
    ],
    resize_keyboard=True
)

kb_genders = ReplyKeyboardMarkup(
    keyboard=[
        [
            KeyboardButton(text="Мужчина"),
            KeyboardButton(text="Женщина"),
        ],
        [
            KeyboardButton(text=BACK_BUTTON)
        ],
    ],
    resize_keyboard=True
)

kb_genders_and_skip = ReplyKeyboardMarkup(
    keyboard=[
        [
            KeyboardButton(text="Мужчинам"),
            KeyboardButton(text="Женщинам"),
        ],
        [
            KeyboardButton(text=SKIP_BUTTON)
        ],
        [
            KeyboardButton(text=BACK_BUTTON)
        ],
    ],
    resize_keyboard=True
)

kb_skip = ReplyKeyboardMarkup(
    keyboard=[
        [
            KeyboardButton(text=SKIP_BUTTON)
        ],
        [
            KeyboardButton(text=BACK_BUTTON)
        ],
    ],
    resize_keyboard=True
)

kb_advertiser = InlineKeyboardMarkup(
    inline_keyboard=[
        [
            InlineKeyboardButton(text="Создать рекламу", callback_data="create_campaign"),
        ],
        [
            InlineKeyboardButton(text="Список реклам", callback_data="all_campaigns"),
        ],
    ],
)

kb_client = InlineKeyboardMarkup(
    inline_keyboard=[
        [
            InlineKeyboardButton(text="Посмотреть рекламу", callback_data="view_ads"),
        ],
    ],
)
