package basic_types

type Subject struct {
	Educators  []string
	Places     []string
	Event_type string
	Event_name string
	Event_time string
}

type Day struct {
	Subjects []Subject
	Date     string
}

const (
	BaseUrl string = "https://mai.ru/education/studies/schedule/"
)

var (
	ShortWeekDays  = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	LongMonthNames = map[string]int{
		"января": 1, "февраля": 2, "марта": 3, "апреля": 4, "мая": 5,
		"июня": 6, "июля": 7, "августа": 8, "сентября": 9, "октября": 10,
		"ноября": 11, "декабря": 12,
	}
)
