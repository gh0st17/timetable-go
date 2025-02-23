package manager

import (
	"github.com/gh0st17/timetable-go/manager/internal/parser"

	"github.com/gh0st17/timetable-go/internal/basic_types"

	"golang.org/x/net/html"
)

func parseSubjects(html_subjects []html.Node, day *basic_types.Day) {
	for i, html_subject := range html_subjects {
		html_subj_name := parser.FindNode(&html_subject, subj_name_param)[0]
		day.Subjects = append(day.Subjects, Subject{})
		parser.ExtractSubject(html_subj_name, &day.Subjects[i])

		html_place := parser.FindNode(&html_subject, place_block_param)[0]
		parser.ExtractPlace(&html_place, &day.Subjects[i])
	}
}

func parseDays(html_days []html.Node, timetable []Day) []Day {
	for _, html_day := range html_days {
		var (
			day           Day
			html_subjects []html.Node
			html_date     *html.Node
		)

		html_date = html_day.FirstChild.NextSibling
		if html_date == nil {
			continue
		}

		day.Date = parser.ExtractText(html_date)

		html_subjects = parser.FindNode(&html_day, subj_param)
		parseSubjects(html_subjects, &day)

		timetable = append(timetable, day)
	}

	return timetable
}
