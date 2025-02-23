// Пакет ical предоставляет набор функции для
// формирования файла ics из среза типа [base_types.Day]
//
// Основные функции:
//   - WriteIcal: Запись расписания в файл ics
//   - CalcWeek: Рассчитывает номер текущей недели
//     с учетом сезона семестра
package ical

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	bt "github.com/gh0st17/timetable-go/manager/internal/basic_types"
	fs "github.com/gh0st17/timetable-go/manager/internal/filesystem"
	"github.com/gh0st17/timetable-go/params"
)

// Запись расписания в файл ics
func WriteIcal(timetable []bt.Day, p *params.Params) error {
	fmt.Printf("Имя файла %s\n", p.FileName)

	var dataString string
	icalDoc := "BEGIN:VCALENDAR\n" + getHeader() + "\n\n\n"

	for _, day := range timetable {
		for i, subject := range day.Subjects {
			dataString = buildDataString(p.GroupName+day.Date, subject)
			uid := stringToHash(dataString)
			icalDoc += getEvent(day, i, uid) + "\n\n"
		}

		icalDoc += "\n"
	}

	icalDoc += "END:VCALENDAR"

	return fs.WriteString(p.OutDir+"/"+p.FileName, icalDoc)
}

// Рассчитывает номер текущей недели с учетом сезона семестра
func CalcWeek() uint {
	today := time.Now()
	_, week := today.ISOWeek()

	if today.Month() >= 8 && today.Day() >= 1 {
		week -= 34
	} else {
		feb8 := time.Date(
			today.Year(), time.February, 8,
			0, 0, 0, 0, today.Location(),
		)
		_, feb8Week := feb8.ISOWeek()
		week -= feb8Week - 1
	}

	if week < 1 {
		week = 1
	} else if week > 18 {
		week = 18
	}

	return uint(week)
}

// Возвращает строку с заголовком ics файла
func getHeader() (header string) {
	return `VERSION:2.0
PRODID:ghost17 | Alexey Sorokin
CALSCALE:GREGORIAN`
}

// Возвращает строки с временем начала и конца занятия
func getDate(month int, day int, hour int, min int) (start string, end string) {
	year := time.Now().Year()

	date := time.Date(
		year, time.Month(month),
		day, hour, min, 0, 0, time.UTC,
	)

	const format string = "20060102T150405"
	start = date.Format(format)
	end = (date.Add(time.Minute * 90)).Format(format)

	return start, end
}

// Возвращает уникальный хэш для события iCal
func stringToHash(dataString string) uint64 {
	hash := sha1.Sum([]byte(dataString))
	return binary.BigEndian.Uint64(hash[:8])
}

// Формирует данные о занятии в строку
func buildDataString(basicString string, subject bt.Subject) string {
	basicString += subject.Event_name + subject.Event_time
	basicString += subject.Event_type

	for _, educator := range subject.Educators {
		basicString += educator
	}

	for _, place := range subject.Places {
		basicString += place
	}

	return basicString
}

// Формирует и возвращает событие iCal для занятия
func getEvent(day bt.Day, eventIdx int, uid uint64) (event string) {
	var (
		startDate string
		endDate   string
		summary   string
		location  string
		subject   bt.Subject = day.Subjects[eventIdx]
		offset    int        = 0
	)

	splittedDate := strings.Split(day.Date, " ")
	if len(splittedDate) < 3 {
		offset = -1
	}

	dayInt, _ := strconv.Atoi(splittedDate[1+offset])
	month := bt.LongMonthNames[splittedDate[2+offset]]

	splittedTime := strings.Split(day.Subjects[eventIdx].Event_time, " ")
	splittedClock := strings.Split(splittedTime[0], ":")
	hour, _ := strconv.Atoi(splittedClock[0])
	min, _ := strconv.Atoi(splittedClock[1])
	startDate, endDate = getDate(month, dayInt, hour, min)

	summary = subject.Event_name

	for _, place := range subject.Places {
		location += fmt.Sprintf("%s / ", place)
	}

	location += subject.Event_type

	for _, educator := range subject.Educators {
		location += fmt.Sprintf(" / %s", educator)
	}

	return fmt.Sprintf(`BEGIN:VEVENT
UID:%d
DTSTART:%s
DTSTAMP:%sZ
DTEND:%s
SUMMARY:%s
LOCATION:%s
END:VEVENT`,
		uid, startDate, startDate, endDate, summary, location)
}
