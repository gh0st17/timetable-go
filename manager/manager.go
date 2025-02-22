package manager

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"timetable/errtype"
	"timetable/internal/basic_types"
	"timetable/internal/database"
	"timetable/manager/parser"
	"timetable/params"

	"golang.org/x/net/html"
)

type Day = basic_types.Day
type Subject = basic_types.Subject
type Params = params.Params

func todayUrl(group *string) string {
	return basic_types.BaseUrl + "index.php?group=" + *group
}

func weekParam(week uint) string {
	return fmt.Sprintf("week=%d", week)
}

func depParam(dep uint) string {
	return fmt.Sprintf("department=Институт+№%d", dep)
}

func courseParam(course uint) string {
	return fmt.Sprintf("course=%d", course)
}

func groupUrl(dep uint, course uint) string {
	return basic_types.BaseUrl + "groups.php?" + depParam(dep) + "&" + courseParam(course)
}

func sessionUrl(group string) string {
	return basic_types.BaseUrl + "session/index.php?group=" + group
}

func fetchGroups(u *url.URL, jar http.CookieJar, proxyUrl *url.URL) ([]string, error) {
	var (
		doc         *html.Node
		err         error
		group_nodes []html.Node
		groups      []string
	)

	pred := func() (*html.Node, error) {
		return loadFromUrl(u, jar, proxyUrl)
	}

	if doc, err = retryLoadFromUrl(3, false, pred); err != nil {
		return nil, err
	}

	parser.FindNode(doc, &group_nodes, &groups_param)

	if len(group_nodes) == 0 {
		return nil, errtype.ErrParse(errors.New("список групп не загружен"))
	}

	for _, group := range group_nodes {
		groups = append(groups, parser.ExtractText(&group))
	}

	sort.Strings(groups)

	return groups, nil
}

func fetchTimetable(doc *html.Node) (timetable []Day, err error) {
	var html_days []html.Node
	parser.FindNode(doc, &html_days, &day_param)

	if len(html_days) == 0 {
		return nil, errtype.ErrParse(errors.New("расписание не найдено"))
	}

	parseDays(&html_days, &timetable)

	return timetable, nil
}

func printTimetable(timetable *[]Day, p *Params) {
	fmt.Printf("Группа %s\n\n", p.GroupName)

	if p.Week != 0 {
		fmt.Printf("Учебная неделя №%d\n\n", p.Week)
	}

	for _, day := range *timetable {
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

func proceedingWeek(p *Params) (u *url.URL) {
	if p.Week != 0 {
		p.FileName += fmt.Sprintf("Week_%d", p.Week)
	}

	if p.Next {
		p.Week = calcWeek()
		p.FileName += fmt.Sprintf("Week_%d", p.Week)
	} else if p.Current {
		p.Week = calcWeek() - 1
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

func Run(p *Params) error {
	var (
		tdb       database.TimetableDB
		doc       *html.Node
		timetable []Day
		u         *url.URL
		err       error
	)

	if err = tdb.LoadDB("timetable.db"); err != nil {
		return err
	}

	if p.WorkDir == "" {
		if p.WorkDir, err = getWd(); err != nil {
			return err
		}

		if p.OutDir == "" {
			p.OutDir = p.WorkDir
		}
	}

	if !dirExists(p.WorkDir + "/groups") {
		createDir(p.WorkDir + "/groups")
	} else if p.Clear {
		removeAllFilesInDir(p.WorkDir + "/groups")
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
	loadCookiesFromFile(jar, "cookies.txt", u)
	if len(jar.Cookies(u)) == 0 {
		_, _ = loadFromUrl(u, jar, p.ProxyUrl)
	}

	pred := func() (*html.Node, error) {
		return loadFromUrl(u, jar, p.ProxyUrl)
	}

	// TO DO
	// Work with timetable in DB at this line

	if doc, err = retryLoadFromUrl(3, true, pred); err != nil {
		return err
	} else {
		// Сохраняем куки в файл
		if err := saveCookiesToFile(jar, "cookies.txt", u); err != nil {
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
		return writeIcal(&timetable, p)
	} else {
		printTimetable(&timetable, p)
	}

	return nil
}
