package manager

import (
	"database/sql"
	"errors"
	"net/http/cookiejar"
	"net/url"
	"timetable/database"
	"timetable/errtype"
)

func proceedingGroupDB(p *Params, tdb *database.TimetableDB, printOnly bool) error {
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
		printLines(&groupsLines, p, printOnly)
	}

	if !printOnly && p.Group == 0 {
		p.GroupName = groupsLines[getUserSelection(&groupsLines)]
	} else if p.Group > 0 && int(p.Group) <= len(groupsLines) {
		p.GroupName = groupsLines[p.Group-1]
	} else if !p.List {
		return errtype.ArgsError(errors.New("номер группы не существует"))
	}

	return nil
}
