package params

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"timetable/errtype"
)

type Params struct {
	Dep       uint8
	Course    uint8
	Week      uint8
	Group     uint8
	GroupName string
	WorkDir   string
	OutDir    string
	FileName  string
	ProxyUrl  *url.URL
	List      bool
	Clear     bool
	Next      bool
	Current   bool
	Session   bool
	Ical      bool
}

func parseUint8(str string) uint8 {
	val, err := strconv.ParseUint(str, 10, 8)
	errtype.HandleError(&err)
	return uint8(val)
}

func printHelp() {
	helpText := `timetable {Институт} {Курс} --group <Число> --week <Число>
timetable {Институт} {Курс} --list
timetable --clear

  Институт      - Номер института от 1 до 12
  Курс          - Номер курса от 1 до 6
  --group,   -g - Номер группы из списка
  --week,    -w - Номер недели от 1 до 18
  --next     -n - Следующая неделя (блокирует -c, -w)
  --current  -c - Текущая неделя (блокирует -w)
  --list,    -l - Показать только список групп
  --ics         - Вывод в ics файл
  --proxy       - Использовать прокси
                  <протокол://адрес:порт>
  --sleep       - Время (в секундах) простоя после загрузки недели для семестра
  --session     - Расписание сессии (блокирует выбор недели: -w, -n, -c)
  --clear       - Очистить кэш групп
  --workdir, -d - Путь рабочей директории (кэш) (по умолчанию равен pwd)
  --output,  -o - Путь для вывода (если не задан то равен -d)
`

	fmt.Println(helpText)
}

func (p *Params) parseArgs(args *[]string) error {
	var (
		err      error
		proxyStr string
		str_ptr  *string
		u8_ptr   *uint8
	)

	for i, arg := range *args {
		if i == 0 {
			continue
		}
		if i == 1 {
			p.Dep = parseUint8(arg)
		} else if i == 2 {
			p.Course = parseUint8(arg)
		} else if arg == "-h" || arg == "--help" {
			printHelp()
			os.Exit(0)
		} else if arg == "-g" || arg == "--group" {
			u8_ptr = &p.Group
		} else if arg == "-w" || arg == "--week" {
			u8_ptr = &p.Week
		} else if arg == "--proxy" {
			str_ptr = &proxyStr
		} else if arg == "--list" {
			p.List = true
		} else if arg == "--clear" {
			p.Clear = true
		} else if arg == "--next" || arg == "-n" {
			p.Next = true
			p.Week = 0
		} else if arg == "--current" || arg == "-c" {
			p.Current = true
			p.Week = 0
		} else if arg == "--session" {
			p.Session = true
		} else if arg == "--ics" {
			p.Ical = true
		} else if arg == "--workdir" || arg == "-d" {
			str_ptr = &p.WorkDir
		} else if arg == "--output" || arg == "-o" {
			str_ptr = &p.OutDir
		} else if u8_ptr == nil && str_ptr == nil {
			printHelp()
			return errtype.ArgsError(fmt.Errorf("неизвестный аргумент '%s'", arg))
		} else if str_ptr != nil {
			*str_ptr = arg
			str_ptr = nil
		} else if res, err := strconv.ParseUint(arg, 10, 8); err == nil {
			*u8_ptr = uint8(res)
			u8_ptr = nil
		} else {
			return errtype.ArgsError(fmt.Errorf("неверное значение '%s', требуется число", arg))
		}
	}

	if proxyStr != "" {
		if p.ProxyUrl, err = url.Parse(proxyStr); err != nil {
			return errtype.RuntimeError(errors.New("неверный формат адреса прокси"))
		}
	}

	return nil
}

func (p *Params) FetchParams() error {
	args := os.Args
	if len(args) < 3 {
		printHelp()
		return errtype.ArgsError(errors.New("недостаточно аргументов"))
	}

	if err := p.parseArgs(&args); err != nil {
		return err
	}

	if p.Week != 0 && p.Week > 18 {
		return errtype.ParseError(errors.New("номер недели должен быть из [1; 18]"))
	}

	return nil
}
