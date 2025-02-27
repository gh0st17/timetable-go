# timetable-go
<p align="center">
  <img width="256" height="256" src="https://github.com/user-attachments/assets/aaa1b413-25a4-4ff3-9577-5487ef99c5f2">
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

- Выбор группы из списка
- Просмотр списка групп
- Загрузка списка групп в кэш
- Загрузка `текущего` расписания, на `конкретную`, `текущую` или `следующую` неделю
- Загрузка расписания `сессии`
- Поддержка `HTTP[S]` и `Socks5` прокси
- Поддержка вывода в формате `iCal`

# Справка по использованию

```
Выбор недели: timetable-go --dep <Институт> --course <Курс> --group <Число> --week <Число>
Список групп: timetable-go --dep <Институт> --course <Курс> --list
Очистка кэша: timetable-go --clear

Флаги:
  -V	Печать номера версии и выход
  -clear
    	Очистить кэш групп
  -course uint
    	Номер курса от 1 до 6
  -current
    	Текущая неделя (игнорирует -week)
  -dep uint
    	Номер института от 1 до 12 или 14
  -group uint
    	Номер группы из списка
  -help
    	Показать эту помощь
  -ics
    	Вывод в ics файл
  -list
    	Показать только список групп
  -next
    	Следующая неделя (игнорирует --current, --week)
  -output string
    	Путь для вывода (если не задан то равен --workdir)
  -proxy string
    	Использовать прокси: <протокол://адрес:порт>
  -session
    	Расписание сессии (игнорирует выбор недели: --week, --next, --current)
  -week uint
    	Номер недели от 1 до 18
  -workdir string
    	Путь рабочей директории (кэш) (по умолчанию равен pwd)
```
