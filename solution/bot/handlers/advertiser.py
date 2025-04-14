import math
from aiogram import Router, F
from aiogram.types import Message, CallbackQuery, ReplyKeyboardRemove, InlineKeyboardMarkup, InlineKeyboardButton
from aiogram.fsm.context import FSMContext
from config import BACK_BUTTON, GRAFANA_URL, SKIP_BUTTON
from handlers.start import cmd_start
from states.advertiser import CreateCampaign
from utils import api
from utils.uuid import generate_uuid_from_id
from keyboards import kb_back, kb_genders_and_skip, kb_advertiser, kb_skip
from states import CreateAdvertiser

advertiser_router = Router()

@advertiser_router.callback_query(F.data == "advertiser_profile")
async def advertiser(call: CallbackQuery, state: FSMContext):
    try:
        await call.message.delete()
    except:
        pass

    advertiser_id = generate_uuid_from_id(call.from_user.id)
    advertiser_profile = await api.get_advertiser_data(advertiser_id)
    if not advertiser_profile:
        await state.set_state(CreateAdvertiser.name)
        await call.message.answer(
            'У Вас еще нет профиля рекламодателя, давайте его создадим.\nКак называется Ваша компания?',
            reply_markup=kb_back,
        )
        return 

    message_text = f"Профиль вашей компании:\nНазвание: {advertiser_profile['name']}"
    await call.message.answer(message_text, reply_markup=kb_advertiser)

@advertiser_router.callback_query(F.data == "all_campaigns")
async def all_campaigns(call: CallbackQuery):
    advertiser_id = generate_uuid_from_id(call.from_user.id)
    count, campaings = await api.get_advertiser_campaigns(advertiser_id, 0)
    if count == 0:
        await call.answer("У вас нет реклам!")
        return
    
    try:
        await call.message.delete()
    except:
        pass


    count_pages = math.ceil(count / 2)
    message_text = f"Ваши рекламные компании:\n\nСтраница 1/{count_pages}"
    buttons = [[]]
    for campaign in campaings:
        buttons[0].append(InlineKeyboardButton(text=campaign['ad_title'][:10], callback_data=f"campaign_open_{campaign['campaign_id']}"))

    if count_pages > 1:
        buttons.append([InlineKeyboardButton(text="➡️", callback_data="page_campaigns_1")])
    
    buttons.append([InlineKeyboardButton(text="Назад", callback_data="advertiser_profile")])

    kb_campaigns = InlineKeyboardMarkup(
        inline_keyboard=buttons
    )

    await call.message.answer(message_text, reply_markup=kb_campaigns)

@advertiser_router.callback_query(F.data.startswith("page_campaigns_"))
async def page_campaigns(call: CallbackQuery):
    page = int(call.data.split("page_campaigns_")[1])

    advertiser_id = generate_uuid_from_id(call.from_user.id)
    count, campaings = await api.get_advertiser_campaigns(advertiser_id, page)

    count_pages = math.ceil(count / 2)
    message_text = f"Ваши рекламные компании:\n\nСтраница {page + 1}/{count_pages}"
    buttons = [[]]
    for campaign in campaings:
        buttons[0].append(InlineKeyboardButton(text=campaign['ad_title'][:10], callback_data=f"campaign_open_{campaign['campaign_id']}"))

    if page != 0:
        buttons.append([InlineKeyboardButton(text="⬅️", callback_data=f"page_campaigns_{page - 1}")])
    if page < count_pages - 1:
        buttons.append([InlineKeyboardButton(text="➡️", callback_data=f"page_campaigns_{page + 1}")])
    
    buttons.append([InlineKeyboardButton(text="Назад", callback_data="advertiser_profile")])

    kb_campaigns = InlineKeyboardMarkup(
        inline_keyboard=buttons
    )
    await call.message.edit_text(message_text, reply_markup=kb_campaigns)

@advertiser_router.callback_query(F.data.startswith("campaign_open_"))
async def page_campaigns(call: CallbackQuery):
    campaign_id = call.data.split("campaign_open_")[1]

    advertiser_id = generate_uuid_from_id(call.from_user.id)
    campaing = await api.get_campaign(advertiser_id, campaign_id)

    message_text = f"{campaing['ad_title']}\n\n{campaing['ad_text']}\n\nЛимит показов: {campaing['impressions_limit']}\nЛимит кликов: {campaing['clicks_limit']}\nЦена за показ: {campaing['cost_per_impression']}\nЦена за клик: {campaing['cost_per_click']}\nДата начала: {campaing['start_date']}\nДата конца: {campaing['end_date']}"
    if campaing['targeting'] != "ALL":
        gender = "Мужской" if campaing['targeting'] == "MALE" else "Женский"
        message_text += f"\nПол для показов: {gender}"
    if campaing.get("age_from"):
        message_text += f"\nОт {campaing.get('age_from')} лет"
    if campaing.get("age_to"):
        message_text += f"\nДо {campaing.get('age_to')} лет"
    if campaing.get("location"):
        message_text += f"\nЛокация: {campaing.get('location')}"
    
    message_text += f'\n\nСтатистика:\n{GRAFANA_URL}/d/fedpld5y7gruob/advertising-platform-dashboard?var-advertiserID={advertiser_id}&var-campaignID={campaign_id}&refresh=30s'

    kb_campaigns = InlineKeyboardMarkup(
        inline_keyboard=[[InlineKeyboardButton(text="В профиль", callback_data="advertiser_profile")]]
    )
    await call.message.edit_text(message_text, reply_markup=kb_campaigns)

@advertiser_router.message(CreateAdvertiser.name, F.text)
async def create_advertiser_name(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    is_created = await api.create_advertiser(generate_uuid_from_id(message.from_user.id), message.text)
    if not is_created:
        await message.answer("Не получилось создать пользователя :(")
    else:
        await message.answer("Профиль рекламодателя успешно создан", keyboard=ReplyKeyboardRemove())
    return await cmd_start(message, state)
    
@advertiser_router.callback_query(F.data == "create_campaign")
async def create_campaign(call: CallbackQuery, state: FSMContext):
    try:
        await call.message.delete()
    except:
        pass

    await call.message.answer("Введите название для рекламного объявления:", reply_markup=kb_back)
    await state.set_state(CreateCampaign.ad_title)

@advertiser_router.message(CreateCampaign.ad_title, F.text)
async def create_advertisement_title(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    ad_title = message.text
    await state.update_data(ad_title=ad_title)
    await message.answer("Введите текст рекламного объявления:")
    await state.set_state(CreateCampaign.ad_text)

@advertiser_router.message(CreateCampaign.ad_text, F.text)
async def create_advertisement_text(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    ad_text = message.text
    await state.update_data(ad_text=ad_text)
    await message.answer("Введите лимит показов:")
    await state.set_state(CreateCampaign.impressions_limit)

@advertiser_router.message(CreateCampaign.impressions_limit, F.text)
async def create_advertisement_impressions_limit(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    try:
        impressions_limit = int(message.text)
        if impressions_limit < 0:
            raise ValueError
    except ValueError:
        await message.answer("Пожалуйста, введите корректный положительный целочисленный лимит показов.")
        return
    await state.update_data(impressions_limit=impressions_limit)
    await message.answer("Введите лимит кликов:")
    await state.set_state(CreateCampaign.clicks_limit)

@advertiser_router.message(CreateCampaign.clicks_limit, F.text)
async def create_advertisement_clicks_limit(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    try:
        clicks_limit = int(message.text)
        impressions_limit = (await state.get_data())["impressions_limit"]
        if clicks_limit < 0 or clicks_limit > impressions_limit:
            raise ValueError
    except ValueError:
        await message.answer("Пожалуйста, введите корректный положительный целочисленный лимит кликов, меньший или равный лимита показов.")
        return
    await state.update_data(clicks_limit=clicks_limit)
    await message.answer("Введите стоимость одного показа:")
    await state.set_state(CreateCampaign.cost_per_impression)

@advertiser_router.message(CreateCampaign.cost_per_impression, F.text)
async def create_advertisement_cost_per_impression(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    try:
        cost_per_impression = float(message.text)
        if cost_per_impression < 0:
            raise ValueError
    except ValueError:
        await message.answer("Пожалуйста, введите корректную положительную цену за показ.")
        return
    await state.update_data(cost_per_impression=cost_per_impression)
    await message.answer("Введите стоимость одного клика:")
    await state.set_state(CreateCampaign.cost_per_click)

@advertiser_router.message(CreateCampaign.cost_per_click, F.text)
async def create_advertisement_cost_per_click(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    try:
        cost_per_click = float(message.text)
        if cost_per_click < 0:
            raise ValueError
    except ValueError:
        await message.answer("Пожалуйста, введите корректную положительную цену за клик.")
        return
        
    await state.update_data(cost_per_click=cost_per_click)
    await message.answer("Введите дату начала показа рекламного объявления:")
    await state.set_state(CreateCampaign.start_date)

@advertiser_router.message(CreateCampaign.start_date, F.text)
async def create_advertisement_start_date(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    current_date = await api.get_current_date()

    try:
        start_date = int(message.text)
        if start_date < 0 or (current_date and start_date < current_date):
            raise ValueError
    except ValueError:
        await message.answer("Пожалуйста, введите корректную дату как число > 0 и не меньше текущей.")
        return
    await state.update_data(start_date=start_date)
    await message.answer("Введите дату окончания показа рекламного объявления:")
    await state.set_state(CreateCampaign.end_date)

@advertiser_router.message(CreateCampaign.end_date, F.text)
async def create_advertisement_end_date(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    current_date = await api.get_current_date()

    try:
        end_date = int(message.text)
        if end_date < 0:
            raise ValueError
        start_date = (await state.get_data())["start_date"]
        if end_date < start_date or (current_date and end_date < current_date):
            raise ValueError
    except ValueError:
        await message.answer("Пожалуйста, введите корректную дату окончания (большую или равную дате начала) и не меньше текущей.")
        return
    await state.update_data(end_date=end_date)
    await message.answer("Введите настройки таргетинга (необязательно). Какому полу показывать вашу рекламу?", reply_markup=kb_genders_and_skip)
    await state.set_state(CreateCampaign.gender)

@advertiser_router.message(CreateCampaign.gender, F.text)
async def create_advertisement_gender(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    if message.text.lower() == 'мужчинам':
        gender = "MALE"
    elif message.text.lower() == "женщинам":
        gender = "FEMALE"
    else:
        gender = "ALL"

    await state.update_data(gender=gender)
    await message.answer("Введите минимальный возраст для таргетинга (необязательно):",reply_markup=kb_skip)
    await state.set_state(CreateCampaign.age_from)

@advertiser_router.message(CreateCampaign.age_from, F.text)
async def create_advertisement_age_from(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    try:
        age_from = int(message.text) if message.text.isdigit() else None
        if age_from and (age_from < 0 or age_from > 100):
            raise ValueError
    except ValueError:
        await message.answer("Пожалуйста, введите корректный минимальный возраст.")
        return
    await state.update_data(age_from=age_from)
    await message.answer("Введите максимальный возраст для таргетинга (необязательно):")
    await state.set_state(CreateCampaign.age_to)

@advertiser_router.message(CreateCampaign.age_to, F.text)
async def create_advertisement_age_to(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    try:
        age_to = int(message.text) if message.text.isdigit() else None
        age_from = (await state.get_data()).get("age_from")
        if (age_to and (age_to < 0 or age_to > 100)) or (age_to and age_from and age_from > age_to):
            raise ValueError
    except ValueError:
        await message.answer("Пожалуйста, введите корректный максимальный возраст.")
        return
    await state.update_data(age_to=age_to)
    await message.answer("Введите локацию для таргетинга (необязательно):")
    await state.set_state(CreateCampaign.location)

@advertiser_router.message(CreateCampaign.location, F.text)
async def create_advertisement_location(message: Message, state: FSMContext):
    if message.text == BACK_BUTTON:
        return await cmd_start(message, state)
    
    location = message.text if message.text != SKIP_BUTTON else None
    await state.update_data(location=location)

    data = await state.get_data()
    impressions_limit = data["impressions_limit"]
    clicks_limit = data["clicks_limit"]
    cost_per_impression = data["cost_per_impression"]
    cost_per_click = data["cost_per_click"]
    ad_title = data["ad_title"]
    ad_text = data["ad_text"]
    start_date = data["start_date"]
    end_date = data["end_date"]
    targeting = {
        "gender": data['gender'],
        "age_from": data.get("age_from"),
        "age_to": data.get("age_to"),
        "location": data.get("location"),
    }

    is_created = await api.create_campaign(
        advertiser_id=generate_uuid_from_id(message.from_user.id),
        impressions_limit=impressions_limit,
        clicks_limit=clicks_limit,
        cost_per_impression=cost_per_impression,
        cost_per_click=cost_per_click,
        ad_title=ad_title,
        ad_text=ad_text,
        start_date=start_date,
        end_date=end_date,
        targeting=targeting,
    )
    
    if is_created:
        await message.answer("Рекламная кампания успешно создана.")
    else:
        await message.answer("Не удалось создать рекламную кампанию.")
    return await cmd_start(message, state)
