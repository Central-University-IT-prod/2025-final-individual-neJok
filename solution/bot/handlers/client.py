from aiogram import Router, F
from aiogram.types import Message, CallbackQuery, ReplyKeyboardRemove, InlineKeyboardButton, InlineKeyboardMarkup
from aiogram.fsm.context import FSMContext
from config import BACK_BUTTON
from handlers.start import cmd_start
from utils import api
from utils.uuid import generate_uuid_from_id
from keyboards import kb_back, kb_genders, kb_client
from states import CreateClient

client_router = Router()

@client_router.callback_query(F.data == "client_profile")
async def client(call: CallbackQuery, state: FSMContext) -> None:
    try:
        await call.message.delete()
    except:
        pass

    client_id = generate_uuid_from_id(call.from_user.id)
    client_profile = await api.get_client_data(client_id)
    if not client_profile:
        await state.set_state(CreateClient.age)
        await call.message.answer(
            'У Вас еще нет профиля клиента, давайте его создадим.\nКак сколько Вам лет?',
            reply_markup=kb_back,
        )
        return

    genders = {
        "MALE": "Мужчина",
        "FEMALE": "Женщина"
    }
    message_text = f"Ваш профиль:\nВозраст: {client_profile['age']}\nЛокация: {client_profile['location']}\nПол: {genders[client_profile['gender']]}"
    await call.message.answer(message_text, reply_markup=kb_client)

@client_router.callback_query(F.data == "view_ads")
async def client(call: CallbackQuery) -> None:
    client_id = generate_uuid_from_id(call.from_user.id)
    ads = await api.get_client_ads(client_id)
    if not ads:
        await call.answer(
            'Реклам пока нет: (',
        )
        return
    
    kb_ad = InlineKeyboardMarkup(
        inline_keyboard=[
            [
                InlineKeyboardButton(text="Перейти", callback_data=f"click_{ads['ad_id']}"),
            ],
        ],
    )

    await call.message.answer(ads['ad_title'] + "\n\n" + ads['ad_text'], reply_markup=kb_ad)

@client_router.callback_query(F.data.startswith("click_"))
async def client(call: CallbackQuery) -> None:
    ad_id = call.data.split("click_")[1]
    client_id = generate_uuid_from_id(call.from_user.id)
    ads = await api.ad_click(client_id, ad_id)
    if ads:
        await call.answer(
            'Успешно!',
        )
    else:
        await call.answer(
            'Ошибка!',
        )
    
@client_router.message(CreateClient.age, F.text)
async def create_client_age(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)

    if not message.text.isdigit() or int(message.text) > 100 or int(message.text) < 0:
        await message.answer("Введите число от 0 до 100!")
        return

    await state.update_data(age=int(message.text))

    await state.set_state(CreateClient.location)
    await message.answer("Напиши свой город проживания:", reply_markup=kb_back)

@client_router.message(CreateClient.location, F.text)
async def create_client_location(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)

    await state.update_data(location=message.text)

    await state.set_state(CreateClient.gender)
    await message.answer("Выбери свой пол:", reply_markup=kb_genders)

@client_router.message(CreateClient.gender, F.text)
async def create_client_location(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    genders = {
        "мужчина": "MALE",
        "женщина": "FEMALE"
    }
    gender = genders.get(message.text.lower())
    if not gender:
        await message.answer("Выберите пол из предложенных ботом :)")
        return


    data = await state.get_data()
    is_created = await api.create_client(generate_uuid_from_id(message.from_user.id), str(message.from_user.username), data['age'], data['location'], gender)
    if not is_created:
        await message.answer("Не получилось создать клиента :(")
        return await cmd_start(message, state)
    
    await message.answer("Профиль клиента успешно создан", keyboard=ReplyKeyboardRemove())
    await cmd_start(message, state)
    