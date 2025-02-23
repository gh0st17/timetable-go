package params

// Строки для справки

const (
	versionDesc = "Печать номера версии и выход"
	versionText = "timetable 1.0.3b2\n" +
		"Copyright (C) 2025\n" +
		"Лицензия MIT: THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY\n" +
		"OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO\n" +
		"THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR\n" +
		"PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR\n" +
		"COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\n" +
		"LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,\n" +
		"ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE\n" +
		"USE OR OTHER DEALINGS IN THE SOFTWARE.\n\n" +
		"Это свободное ПО: вы можете изменять и распространять его.\n" +
		"Нет НИКАКИХ ГАРАНТИЙ в пределах действующего законодательства.\n\n" +
		"Автор: Alexey Sorokin.\n"

	weekExample  = "--dep <Институт> --course <Курс> --group <Число> --week <Число>"
	listExample  = "--dep <Институт> --course <Курс> --list"
	clearExample = "--clear"

	depDesc         = "Номер института от 1 до 12 или 14"
	courseDesc      = "Номер курса от 1 до 6"
	groupDesc       = "Номер группы из списка"
	weekDesc        = "Номер недели от 1 до 18"
	nextWeekDesc    = "Следующая неделя (игнорирует --current, --week)"
	currentWeekDesc = "Текущая неделя (игнорирует -week)"
	listDesc        = "Показать только список групп"
	icsDesc         = "Вывод в ics файл"
	helpDesc        = "Показать эту помощь"
	proxyDesc       = "Использовать прокси: <протокол://адрес:порт>"
	sessionDesc     = "Расписание сессии (игнорирует выбор недели: --week, --next, --current)"
	clearDesc       = "Очистить кэш групп"
	workdirDesc     = "Путь рабочей директории (кэш) (по умолчанию равен pwd)"
	outPathDesc     = "Путь для вывода (если не задан то равен --workdir)"
)
