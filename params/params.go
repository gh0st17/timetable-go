// Пакет params предоставляет набор функции для
// обработки входных параметров программы
//
// Основные функции:
//   - ParseParams: Обрабатывает входные флаги и возвращает
//     структуру [Params] с результатом обработанных флагов
package params

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gh0st17/timetable-go/errtype"
)

type Params struct {
	Dep       uint
	Course    uint
	Week      uint
	Group     uint
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

// Печатает справку
func printHelp() {
	program := filepath.Base(os.Args[0])

	fmt.Println("Выбор недели:", program, weekExample)
	fmt.Println("Список групп:", program, listExample)
	fmt.Println("Очистка кэша:", program, clearExample)
	fmt.Printf("\nФлаги:\n")

	flag.PrintDefaults()
}

// Возвращает структуру Params с прочитанными
// входными аргументами программы
func ParseParams() (p *Params, err error) {
	p = &Params{}
	flag.Usage = printHelp
	flag.StringVar(&p.OutDir, "output", "", outPathDesc)
	flag.StringVar(&p.WorkDir, "workdir", "", workdirDesc)

	var proxyStr string
	flag.StringVar(&proxyStr, "proxy", "", proxyDesc)

	flag.UintVar(&p.Dep, "dep", 0, depDesc)
	flag.UintVar(&p.Course, "course", 0, courseDesc)
	flag.UintVar(&p.Group, "group", 0, groupDesc)
	flag.UintVar(&p.Week, "week", 0, weekDesc)

	flag.BoolVar(&p.List, "list", false, listDesc)
	flag.BoolVar(&p.Clear, "clear", false, clearDesc)
	flag.BoolVar(&p.Next, "next", false, nextWeekDesc)
	flag.BoolVar(&p.Current, "current", false, currentWeekDesc)
	flag.BoolVar(&p.Session, "session", false, sessionDesc)
	flag.BoolVar(&p.Ical, "ics", false, icsDesc)

	version := flag.Bool("V", false, versionDesc)
	help := flag.Bool("help", false, helpDesc)

	flag.Parse()

	if *version {
		fmt.Print(versionText)
		os.Exit(0)
	}
	if *help {
		printHelp()
		os.Exit(0)
	}
	if p.Clear {
		return p, nil
	}

	if err = checkDep(p); err != nil {
		return nil, errtype.ErrArgument(err)
	} else if err = checkCourse(p); err != nil {
		return nil, errtype.ErrArgument(err)
	} else if err = checkWeek(p); err != nil {
		return nil, errtype.ErrArgument(err)
	}

	return p, nil
}

// Проверяет корректность ввода номера института
func checkDep(p *Params) error {
	if p.Dep < 1 || p.Dep > 14 || p.Dep == 13 {
		return ErrDepOutOfBound
	}
	return nil
}

// Проверяет корректность ввода номера курса
func checkCourse(p *Params) error {
	if p.Course < 1 || p.Course > 6 {
		return ErrCourseOutOfBound
	}

	return nil
}

// Проверяет корректность ввода номера недели
func checkWeek(p *Params) error {
	if p.Session {
		p.Week = 0
		p.Current = false
		p.Next = false
	} else if p.Next {
		p.Week = 0
		p.Current = false
	} else if p.Current {
		p.Week = 0
	}

	if p.Week != 0 && p.Week > 18 {
		return ErrWeekInput
	}

	return nil
}
