// Пакет manager предоставляет набор функции для
// организации работы всех функции программы в
// в зависимости от входных параметров [params.Params]
//
// Основные функции:
//   - Run: Запускает работу программы
package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"

	"github.com/gh0st17/timetable-go/errtype"
	bt "github.com/gh0st17/timetable-go/manager/internal/basic_types"
	db "github.com/gh0st17/timetable-go/manager/internal/database"
	fs "github.com/gh0st17/timetable-go/manager/internal/filesystem"
	"github.com/gh0st17/timetable-go/manager/internal/ical"
	"github.com/gh0st17/timetable-go/manager/internal/parser"
	"github.com/gh0st17/timetable-go/manager/internal/userio"
	"github.com/gh0st17/timetable-go/params"

	"golang.org/x/net/html"
)

// Запускает работу программы
func Run(p *params.Params) error {
	var (
		tdb       db.TimetableDB
		doc       *html.Node
		timetable []bt.Day
		u         *url.URL
		err       error
	)

	if err = tdb.LoadDB("timetable.db"); err != nil {
		return err
	}

	if p.WorkDir == "" {
		if p.WorkDir, err = fs.GetWd(); err != nil {
			return errtype.ErrRuntime(err)
		}

		if p.OutDir == "" {
			p.OutDir = p.WorkDir
		}
	}

	if p.Clear {
		return tdb.Delete("groups", []db.Criteria{})
	}

	if err = proceedingGroupDB(p, &tdb, p.List); err != nil {
		return err
	}

	p.FileName = p.GroupName + "_"

	if p.List {
		return nil
	}

	if p.Session {
		p.FileName += "Session.ics"
		u, _ = url.Parse(sessionUrl(p.GroupName))
	} else {
		u = proceedingWeek(p)
	}

	jar, _ := cookiejar.New(nil)
	if err = fs.LoadCookiesFromFile(jar, "cookies.txt", u); err != nil {
		return errtype.ErrRuntime(err)
	}
	if len(jar.Cookies(u)) == 0 {
		_, _ = loadFromUrl(u, jar, p.ProxyUrl)
	}

	pred := func() (*html.Node, error) {
		return loadFromUrl(u, jar, p.ProxyUrl)
	}

	// TO DO
	// Work with timetable in DB at this line

	if doc, err = retryLoadFromUrl(3, true, pred); err != nil {
		return errtype.ErrNetwork(err)
	} else {
		// Сохраняем куки в файл
		if err := fs.SaveCookiesToFile(jar, "cookies.txt", u); err != nil {
			return errtype.ErrRuntime(fmt.Errorf("ошибка сохранения куки: %s", err))
		}
	}

	if timetable, err = fetchTimetable(doc); err != nil {
		return err
	}

	if err = tdb.CloseDB(); err != nil {
		return err
	}

	if p.Ical {
		return ical.WriteIcal(timetable, p)
	} else {
		printTimetable(timetable, p)
	}

	return nil
}

// Возвращает ссылку текущего расписания
func todayUrl(group *string) string {
	return bt.BaseUrl + "index.php?group=" + *group
}

// Возвращает параметр номера недели
func weekParam(week uint) string {
	return fmt.Sprintf("week=%d", week)
}

// Возвращает параметр номера института
func depParam(dep uint) string {
	return fmt.Sprintf("department=Институт+№%d", dep)
}

// Возвращает параметр номера курса
func courseParam(course uint) string {
	return fmt.Sprintf("course=%d", course)
}

// Возвращает адрес страницы с выбором группы
func groupUrl(dep uint, course uint) string {
	return bt.BaseUrl + "groups.php?" + depParam(dep) + "&" + courseParam(course)
}

// Возвращает адрес страницы с расписанием сессии
func sessionUrl(group string) string {
	return bt.BaseUrl + "session/index.php?group=" + group
}

// Загружает список групп по сети
func fetchGroups(u *url.URL, jar http.CookieJar, proxyUrl *url.URL) ([]string, error) {
	var (
		doc         *html.Node
		err         error
		group_nodes = []html.Node{}
		groups      []string
	)

	pred := func() (*html.Node, error) {
		return loadFromUrl(u, jar, proxyUrl)
	}

	if doc, err = retryLoadFromUrl(3, false, pred); err != nil {
		return nil, err
	}

	group_nodes = parser.FindNode(doc, groups_param)

	if len(group_nodes) == 0 {
		return nil, errtype.ErrParse(errors.New("список групп не загружен"))
	}

	for _, group := range group_nodes {
		groups = append(groups, parser.ExtractText(&group))
	}

	sort.Strings(groups)

	return groups, nil
}

// Выполняет разбор страницы с расписанием
func fetchTimetable(doc *html.Node) (timetable []bt.Day, err error) {
	html_days := parser.FindNode(doc, day_param)

	if len(html_days) == 0 {
		return nil, errtype.ErrParse(errors.New("расписание не найдено"))
	}

	timetable = parseDays(html_days, timetable)

	return timetable, nil
}

// Печатает расписание в окно консоли
func printTimetable(timetable []bt.Day, p *params.Params) {
	fmt.Printf("Группа %s\n\n", p.GroupName)

	if p.Week != 0 {
		fmt.Printf("Учебная неделя №%d\n\n", p.Week)
	}

	for _, day := range timetable {
		fmt.Println(day.Date)
		for _, subject := range day.Subjects {
			fmt.Printf(
				"[%s] %s\n%s",
				subject.Event_type, subject.Event_name, subject.Event_time,
			)

			for _, educator := range subject.Educators {
				fmt.Printf(" / %s", educator)
			}

			for _, place := range subject.Places {
				fmt.Printf(" / %s", place)
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

// Обработка части имени файла ics
func proceedingWeek(p *params.Params) (u *url.URL) {
	if p.Week != 0 {
		p.FileName += fmt.Sprintf("Week_%d", p.Week)
	}

	if p.Next {
		p.Week = ical.CalcWeek()
		p.FileName += fmt.Sprintf("Week_%d", p.Week)
	} else if p.Current {
		p.Week = ical.CalcWeek() - 1
		p.FileName += fmt.Sprintf("Week_%d", p.Week)
	} else if p.Week == 0 {
		u, _ = url.Parse(todayUrl(&p.GroupName))
		p.FileName += "Today.ics"
		return u
	}

	p.FileName += ".ics"

	u, _ = url.Parse(todayUrl(&p.GroupName) + "&" + weekParam(p.Week))
	return u
}

// Записывает название группы в p, прочитанное из базы данных
// или из пользовательского ввода
func proceedingGroupDB(p *params.Params, tdb *db.TimetableDB, printOnly bool) error {
	var (
		groupsLines []string
		rows        *sql.Rows
		err         error
	)

	if rows, err = tdb.QueryGroup(p.Dep, p.Course); err != nil {
		return err
	}
	defer rows.Close()

	if groupsLines, err = tdb.GetGroupsLines(rows); err != nil {
		return err
	}

	if len(groupsLines) == 0 {
		u, _ := url.Parse(groupUrl(p.Dep, p.Course))
		jar, _ := cookiejar.New(nil)

		if groupsLines, err = fetchGroups(u, jar, p.ProxyUrl); err != nil {
			return err
		}

		if err = tdb.InsertGroup(groupsLines, p); err != nil {
			return err
		}
	}

	if p.Group == 0 {
		userio.PrintGroupLines(groupsLines, p, printOnly)
	}

	if !printOnly && p.Group == 0 {
		p.GroupName = groupsLines[userio.GetUserSelection(groupsLines)]
	} else if p.Group > 0 && int(p.Group) <= len(groupsLines) {
		p.GroupName = groupsLines[p.Group-1]
	} else if !p.List {
		return errtype.ErrArgument(errors.New("номер группы не существует"))
	}

	return nil
}
