package main

import (
	"timetable/errtype"
	"timetable/manager"
	"timetable/params"
)

func main() {
	var p = params.Params{}
	if err := p.FetchParams(); err != nil {
		errtype.HandleError(&err)
	}

	if err := manager.Run(&p); err != nil {
		errtype.HandleError(&err)
	}
}
