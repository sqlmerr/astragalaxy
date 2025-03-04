main_menu = 📋 <b><i>главное меню</i></b>
spaceship_menu_btn = 🚀 корабли
spaceship_menu_choose_spaceship = 🚀 <b><i>Выберите корабль:</i></b>
spaceship_menu =
    <b><i>🚀 твой корабль</i></b>:
      🪪 <i>имя:</i> <code>{ $name }</code>
spaceship_menu_exit = ⛔️ выйти
spaceship_menu_enter = ❇️ войти
spaceship_menu_change_name = 💠 переименовать
spaceship_menu_enter_name = 💠 <b><i>отправьте новое имя</i></b>
    <blockquote>🪪 текущее имя: { $name }</blockquote>
invalid_spaceship_name = ❌ <b><i>имя не подходит по критериям:</i></b>
      - длина больше 3, но меньше 32
      - не содержит пробелов
      - в ней содержатся только русские буквы, латинские буквы, цифры, @, #, $, %, !, *, _, -
      - не может быть только цифрам🌃и
spaceship_menu_copy_name = скопировать имя

starmap_menu = 🌃 <b><i>звездная карта</i></b>
starmap_system_info = ☀️ <b>система</b> <code>{ $name }</code>
    - <b>id</b>: <code>{ $id }</code>

starmap_planet = 🪐 <b>планета</b> <code>{ $name }</code>
    - <b>находится в системе</b> <code>{ $system_name }</code>
    <code>............</code>
    ⚠️ <b>опасность:</b> { $threat ->
        *[other] none
        [toxins] токсичность
        [radiation] радиация
        [heat] жара
        [freezing] холод
    }
starmap-flight-info = 🛸 <b>летим</b>
    - <i>осталось времени:</i> { $time } с.


flight-error = ошибка
flight-error-already_flying = Вы уже летите!
flight-success = ✅
flight-not_in_spaceship = ⚠ Ты должен сидеть в корабле!