package manager

import "timetable/manager/parser"

// Параметры поиска тэгов для страницы расписания
var (
	day_param = parser.NodeParam{
		Tag:       "div",
		Attr_name: "class",
		Attr_val:  "step-content",
	}
	subj_param = parser.NodeParam{
		Tag:       "div",
		Attr_name: "class",
		Attr_val:  "mb-4",
	}
	subj_name_param = parser.NodeParam{
		Tag:       "p",
		Attr_name: "class",
		Attr_val:  "mb-2 fw-semi-bold text-dark",
	}
	place_block_param = parser.NodeParam{
		Tag:       "ul",
		Attr_name: "class",
		Attr_val:  "list-inline list-separator text-body small",
	}
)

// Парамтеры поиска тэгов для страницы с группами
var (
	groups_param = parser.NodeParam{
		Tag:       "a",
		Attr_name: "class",
		Attr_val:  "btn btn-soft-secondary btn-xs mb-1 fw-medium btn-group",
	}
)
