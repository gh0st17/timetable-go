# TimeTable
<p align="center">
  <img width="256" height="256" src="https://i.ibb.co/FXzfktM/Default-logo.png">
</p>

<p align="center">
  Парсер расписания с сайта МАИ
</p>


<p align="center">
  <a href="https://github.com/gh0st17/timetable-go/releases/latest"><img src="https://img.shields.io/github/v/release/gh0st17/timetable-go?style=plastic"></a>
  <img src="https://img.shields.io/badge/license-MIT-blue?style=plastic">
  <img src="https://tokei.rs/b1/github/gh0st17/timetable-go?category=code">
</p>

# Возможности:

- [x] Выбор группы из списка
- [x] Просмотр списка групп
- [x] Загрузка списка групп в кэш
- [ ] Загрузка `текущего` рассписания, на `конкретную`, `текущую` или `следующую` неделю
- [ ] Загрузка рассписания `сессии`
- [x] Поддержка `HTTP[S]` и `Socks5` прокси
- [ ] Поддержка вывода в формате `iCal`

# Справка по использованию

```
timetable {Институт} {Курс} --group <Число> --week <Число>
timetable {Институт} {Курс} --list
timetable --clear

  Институт      - Номер института от 1 до 12
  Курс          - Номер курса от 1 до 6
  --group,   -g - Номер группы из списка
  --week,    -w - Номер недели от 1 до 18 или current для текущей недили, next — для следующей
  --list,    -l - Показать только список групп
  --ics         - Вывод в ics файл
  --proxy       - Использовать прокси
                  <протокол://адрес:порт>
  --sleep       - Время (в секундах) простоя после загрузки недели для семестра
  --session     - Расписание сессии
  --clear       - Очистить кэш групп
  --workdir, -d - Путь рабочей директории (кэш)
  --output,  -o - Путь для вывода
```
