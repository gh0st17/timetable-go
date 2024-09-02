package manager

import (
	"timetable/basic_types"
	"timetable/manager/parser"

	"golang.org/x/net/html"
)

func parseSubjects(html_subjects *[]html.Node, day *basic_types.Day) {
	var (
		html_subj_name []html.Node
		html_place     []html.Node
	)

	for i, html_subject := range *html_subjects {
		parser.FindNode(&html_subject, &html_subj_name, &subj_name_param)
		if html_subj_name != nil {
			day.Subjects = append(day.Subjects, Subject{})
			parser.ExtractSubject(&html_subj_name, &day.Subjects[i])
		}

		parser.FindNode(&html_subject, &html_place, &place_block_param)
		parser.ExtractPlace(&html_place[i], &day.Subjects[i])
	}
}

func parseDays(html_days *[]html.Node, timetable *[]Day) {
	for _, html_day := range *html_days {
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

		parser.FindNode(&html_day, &html_subjects, &subj_param)
		parseSubjects(&html_subjects, &day)

		*timetable = append(*timetable, day)
	}
}
