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
	LongMonthNames = []string{
		"января", "февраля", "марта", "апреля", "мая",
		"июня", "июля", "августа", "сентября", "октября",
		"ноября", "декабря",
	}
)
