package database

import (
	"database/sql"
	"fmt"
	"time"
	"timetable/errtype"
	"timetable/params"

	_ "github.com/mattn/go-sqlite3"
)

type Group struct {
	Year         uint16
	SemesterTime uint8
	Department   uint8
	Course       uint8
	GroupName    string
}

type Subject struct {
	Educator1  string
	Educator2  string
	Place1     string
	Place2     string
	Event_type string
	Event_name string
	Event_time string
}

type TimetableDB struct {
	tdb *sql.DB
}

func CalculateSemType() (semesterTime uint8) {
	now := time.Now()
	if now.Month() > 8 {
		semesterTime = 1
	} else {
		semesterTime = 0
	}

	return semesterTime
}

func buildGroup(p *params.Params) *Group {
	return &Group{
		Year:         uint16(time.Now().Year()),
		SemesterTime: CalculateSemType(),
		Department:   p.Dep,
		Course:       p.Course,
		GroupName:    "",
	}
}

func (db *TimetableDB) LoadDB(fileName string) error {
	var err error
	db.tdb, err = sql.Open("sqlite3", fileName)
	if err != nil {
		return errtype.RuntimeError(
			fmt.Errorf("ошибка открытия базы данных: %s", err),
		)
	}

	return nil
}

func (db *TimetableDB) CloseDB() error {
	if err := db.tdb.Close(); err != nil {
		return errtype.RuntimeError(
			fmt.Errorf("ошибка закрытия базы данных: %s", err),
		)
	}

	return nil
}

func (db *TimetableDB) InsertGroup(groupsLines []string, p *params.Params) error {
	var (
		err   error
		query string
	)

	query = `INSERT INTO groups 
(year, semesterTime, department, course, groupName) VALUES
`

	groupQ := buildGroup(p)
	for i, group := range groupsLines {
		query += fmt.Sprintf("(%d, %d, %d, %d, '%s')",
			groupQ.Year, groupQ.SemesterTime, groupQ.Department,
			groupQ.Course, group)

		if i+1 != len(groupsLines) {
			query += ",\n"
		}
	}

	if _, err = db.tdb.Exec(query); err != nil {
		return errtype.RuntimeError(
			fmt.Errorf("ошибка при добавлении записи в базу данных: %s", err),
		)
	}

	return nil
}

func (db *TimetableDB) QueryGroup(dep uint8, course uint8) (*sql.Rows, error) {
	var (
		err   error
		query string
		rows  *sql.Rows
	)

	query = fmt.Sprintf(`SELECT groupName FROM groups 
WHERE department=%d AND course=%d
ORDER BY groupName ASC`,
		dep, course)

	if rows, err = db.tdb.Query(query); err != nil {
		return nil, errtype.RuntimeError(
			fmt.Errorf("ошибка при запросе групп в базе данных: %s", err),
		)
	}

	return rows, nil
}

func (db *TimetableDB) GetGroupsLines(rows *sql.Rows) ([]string, error) {
	var (
		line        string
		groupsLines []string
	)

	for rows.Next() {
		err := rows.Scan(&line)
		if err != nil {
			return nil, errtype.DatabaseError(
				fmt.Errorf("ошибка чтения базы данных: %s", err),
			)
		}
		groupsLines = append(groupsLines, line)
	}

	return groupsLines, nil
}

func (db *TimetableDB) DeleteGroups(dep uint8, course uint8, semType uint8) error {
	query := fmt.Sprintf(`DELETE FROM groups 
WHERE department=%d AND course=%d AND semesterType=%d`,
		dep, course, semType)

	if _, err := db.tdb.Exec(query); err != nil {
		return errtype.RuntimeError(
			fmt.Errorf("ошибка при запросе удаления групп в базе данных: %s", err),
		)
	}

	return nil
}
