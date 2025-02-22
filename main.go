package main

import (
	"github.com/gh0st17/timetable-go/errtype"
	"github.com/gh0st17/timetable-go/manager"
	"github.com/gh0st17/timetable-go/params"
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
