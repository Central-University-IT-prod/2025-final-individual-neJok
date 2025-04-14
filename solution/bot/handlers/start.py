from aiogram import Router, F
from aiogram.filters import CommandStart
from aiogram.types import Message, ReplyKeyboardRemove
from aiogram.fsm.context import FSMContext
from keyboards import kb_choose_profile

start_router = Router()

@start_router.message(CommandStart(), F.chat.type == "private")
async def cmd_start(message: Message, state: FSMContext):
    await state.clear()
    msg = await message.answer("Убираю кнопки...", reply_markup=ReplyKeyboardRemove())
    await msg.delete()
    await message.answer(
        f'👋 Привет, {message.from_user.first_name}!\n\n'
        'В какой личный кабинет Вы хотите зайти? 👇',
        reply_markup=kb_choose_profile,
    )
