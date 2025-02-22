package main

import (
	"timetable/errtype"
	"timetable/manager"
	"timetable/params"
)

func main() {
	if p, err := params.ParseParams(); err != nil {
		errtype.ErrorHandler(errtype.ErrArgument(err))
	} else {
		if err := manager.Run(p); err != nil {
			errtype.ErrorHandler(err)
		}
	}
}
