package manager

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"timetable/errtype"
	"timetable/manager/basic_types"
	"timetable/manager/parser"
	"timetable/params"

	"golang.org/x/net/html"
)

// import (
// 	"time"
// )

func todayUrl(group *string) string {
	return basic_types.BaseUrl + "index.php?group=" + *group
}

func weekParam(week uint8) string {
	return fmt.Sprintf("week=%d", week)
}

func depParam(dep uint8) string {
	return fmt.Sprintf("department=Институт+№%d", dep)
}

func courseParam(course uint8) string {
	return fmt.Sprintf("course=%d", course)
}

func groupUrl(dep uint8, course uint8) string {
	return basic_types.BaseUrl + "groups.php?" + depParam(dep) + "&" + courseParam(course)
}

// func sessionUrl(group string) string {
// 	return basic_types.BaseUrl + "session/index.php?group=" + group
// }

func fetchGroups(u *url.URL, jar http.CookieJar, proxyUrl *url.URL) ([]string, error) {
	var (
		doc         *html.Node
		err         error
		group_nodes []html.Node
		groups      []string
	)

	if doc, err = load_from_url(u, jar, proxyUrl); err != nil {
		return nil, err
	}

	parser.FindNode(doc, &group_nodes, &groups_param)

	if len(group_nodes) == 0 {
		return nil, errtype.ParseError(errors.New("список групп не загружен"))
	}

	for _, group := range group_nodes {
		groups = append(groups, parser.ExtractText(&group))
	}

	sort.Strings(groups)

	return groups, nil
}

func fetchTimetable(doc *html.Node) (timetable []basic_types.Day, err error) {
	var html_days []html.Node
	parser.FindNode(doc, &html_days, &day_param)

	if len(html_days) == 0 {
		return nil, errtype.ParseError(errors.New("расписание не найдено"))
	}

	parseDays(&html_days, &timetable)

	return timetable, nil
}

func printTimetable(timetable *[]basic_types.Day, p *params.Params) {
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

func proceedingEmptyGroup(p *params.Params, printOnly bool) error {
	u, _ := url.Parse(groupUrl(p.Dep, p.Course))
	jar, _ := cookiejar.New(nil)
	groupFile := fmt.Sprintf("%s/groups/%d-%d.txt", p.WorkDir, p.Dep, p.Course)

	var (
		groups []string
		err    error
	)

	if fileExists(groupFile) {
		if groups, err = readFile(groupFile); err != nil {
			return err
		}
	} else {
		if groups, err = fetchGroups(u, jar, p.ProxyUrl); err != nil {
			return err
		}
		if err = writeFile(groupFile, &groups); err != nil {
			return err
		}
	}

	if p.Group == 0 {
		printLines(&groups, p, printOnly)
	}

	if !printOnly && p.Group == 0 {
		p.GroupName = groups[getUserSelection(&groups)]
	} else if p.Group > 0 && int(p.Group) <= len(groups) {
		p.GroupName = groups[p.Group-1]
	} else {
		return errtype.ArgsError(errors.New("номер группы не существует"))
	}

	return nil
}

func Run(p *params.Params) error {
	var (
		u   *url.URL
		err error
	)
	jar, _ := cookiejar.New(nil)

	if p.WorkDir == "" {
		if p.WorkDir, err = getWd(); err != nil {
			return err
		}
	}

	if !dirExists(p.WorkDir + "/groups") {
		createDir(p.WorkDir + "/groups")
	}

	if err = proceedingEmptyGroup(p, p.List); err != nil {
		return err
	}

	if p.List {
		return nil
	}

	if p.Week == 0 {
		u, _ = url.Parse(todayUrl(&p.GroupName))
	} else {
		u, _ = url.Parse(todayUrl(&p.GroupName) + "&" + weekParam(p.Week))
	}

	loadCookiesFromFile(jar, "cookies.txt", u)
	if len(jar.Cookies(u)) == 0 {
		load_from_url(u, jar, p.ProxyUrl)
	}

	if doc, err := load_from_url(u, jar, p.ProxyUrl); err == nil {
		if timetable, err := fetchTimetable(doc); err == nil {
			printTimetable(&timetable, p)
		} else {
			return err
		}
	} else {
		return err
	}

	// Сохраняем куки в файл
	if err := saveCookiesToFile(jar, "cookies.txt", u); err != nil {
		return errtype.RuntimeError(fmt.Errorf("ошибка сохранения куки: %s", err))
	}

	return nil
}
