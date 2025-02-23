package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gh0st17/timetable-go/errtype"
	"github.com/gh0st17/timetable-go/params"

	_ "github.com/mattn/go-sqlite3"
)

// Прдеставляет таблицу groups
type group struct {
	Year         uint16
	SemesterTime uint8
	Department   uint
	Course       uint
	GroupName    string
}

// Прдеставляет таблицу timetable
type subject struct {
	Educator1  string
	Educator2  string
	Place1     string
	Place2     string
	Event_type string
	Event_name string
	Event_time string
}

// Представляет критрии для подстановки в условие
// SQL запроса
type Criteria struct {
	Key          string
	Value        any
	PostOperator string
}

type TimetableDB struct {
	tdb *sql.DB
}

// Возвращает 0 если семестр - весенний, иначе возващает 1
func CalculateSemType() (semesterTime uint8) {
	now := time.Now()
	if now.Month() > 8 {
		semesterTime = 1
	} else {
		semesterTime = 0
	}

	return semesterTime
}

// Возвращает запись группы
func buildGroup(p *params.Params) *group {
	return &group{
		Year:         uint16(time.Now().Year()),
		SemesterTime: CalculateSemType(),
		Department:   p.Dep,
		Course:       p.Course,
		GroupName:    "",
	}
}

// Загружает локальную базу данных из файла
func (db *TimetableDB) LoadDB(fileName string) error {
	var err error
	db.tdb, err = sql.Open("sqlite3", fileName)
	if err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrOpenDB, err))
	}

	return nil
}

// Закрывает базу данных
func (db *TimetableDB) CloseDB() error {
	if err := db.tdb.Close(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrCloseDB, err))
	}

	return nil
}

// Вставка в таблицу groups
func (db *TimetableDB) InsertGroup(groupsLines []string, p *params.Params) error {
	var (
		err   error
		query string
	)

	query = "INSERT INTO groups " +
		"(year, semesterTime, department, course, groupName) VALUES "

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
		return errtype.ErrDataBase(errtype.Join(ErrInsetGroup, err))
	}

	return nil
}

// Запрос в таблице groups
func (db *TimetableDB) QueryGroup(dep uint, course uint) (*sql.Rows, error) {
	criteries := []Criteria{}
	c := Criteria{
		Key:          "department",
		Value:        dep,
		PostOperator: "AND",
	}
	criteries = append(criteries, c)

	c = Criteria{
		Key:          "course",
		Value:        course,
		PostOperator: "",
	}
	criteries = append(criteries, c)

	return db.query("groupName", "groups", criteries)
}

// Общая функция для запросов в базе данных
func (db *TimetableDB) query(sel string, table string, criteries []Criteria) (*sql.Rows, error) {
	var (
		err   error
		query string
		rows  *sql.Rows
	)

	query = fmt.Sprintf("SELECT %s FROM %s WHERE ", sel, table)
	for _, c := range criteries {
		query += fmt.Sprintf("%s=%v %s ", c.Key, c.Value, c.PostOperator)
	}
	query += fmt.Sprintf("ORDER BY %s ASC", sel)

	if rows, err = db.tdb.Query(query); err != nil {
		return nil, errtype.ErrDataBase(errtype.Join(ErrQueryGroup, err))
	}

	return rows, nil
}

// Возвраает список групп в виде среза строк
func (db *TimetableDB) GetGroupsLines(rows *sql.Rows) ([]string, error) {
	var (
		line        string
		groupsLines []string
	)

	for rows.Next() {
		err := rows.Scan(&line)
		if err != nil {
			return nil, errtype.ErrDataBase(errtype.Join(ErrReadGroups, err))
		}
		groupsLines = append(groupsLines, line)
	}

	return groupsLines, nil
}

// Общий метод для удаления записей из таблиц базы данных
func (db *TimetableDB) Delete(table string, criteries []Criteria) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE ", table)

	if len(criteries) > 0 {
		for _, c := range criteries {
			query += fmt.Sprintf("%s=%v %s ", c.Key, c.Value, c.PostOperator)
		}
	} else {
		query += "1;"
	}

	if res, err := db.tdb.Exec(query); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrDelete, err))
	} else {
		affected, _ := res.RowsAffected()
		fmt.Println("Удалено строк:", affected)
	}

	return nil
}
