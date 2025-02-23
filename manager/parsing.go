package manager

import (
	bt "github.com/gh0st17/timetable-go/manager/internal/basic_types"
	"github.com/gh0st17/timetable-go/manager/internal/parser"

	"golang.org/x/net/html"
)

// Разбор предметов
func parseSubjects(html_subjects []html.Node, day bt.Day) bt.Day {
	for i, html_subject := range html_subjects {
		html_subj_name := parser.FindNode(&html_subject, subj_name_param)[0]
		day.Subjects = append(day.Subjects, bt.Subject{})
		e_type, e_name := parser.ExtractSubject(html_subj_name)

		html_place := parser.FindNode(&html_subject, place_block_param)[0]
		day.Subjects[i] = parser.ExtractPlace(&html_place)
		day.Subjects[i].Event_name = e_name
		day.Subjects[i].Event_type = e_type
	}

	return day
}

// Разбор учебных дней
func parseDays(html_days []html.Node, timetable []bt.Day) []bt.Day {
	for _, html_day := range html_days {
		var (
			day           bt.Day
			html_subjects []html.Node
			html_date     *html.Node
		)

		html_date = html_day.FirstChild.NextSibling
		if html_date == nil {
			continue
		}

		day.Date = parser.ExtractText(html_date)

		html_subjects = parser.FindNode(&html_day, subj_param)
		day = parseSubjects(html_subjects, day)

		timetable = append(timetable, day)
	}

	return timetable
}
