package errtype

import (
	"fmt"
	"os"
)

type Error struct {
	Err  error
	Code int
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func ArgsError(err error) *Error {
	return &Error{
		Err:  err,
		Code: 1,
	}
}

func NetworkError(err error) *Error {
	return &Error{
		Err:  err,
		Code: 2,
	}
}

func ParseError(err error) *Error {
	return &Error{
		Err:  err,
		Code: 3,
	}
}

func RuntimeError(err error) *Error {
	return &Error{
		Err:  err,
		Code: 4,
	}
}

func DatabaseError(err error) *Error {
	return &Error{
		Err:  err,
		Code: 5,
	}
}

func HandleError(err *error) {
	if *err != nil {
		fmt.Println(*err)
		if e, ok := (*err).(*Error); ok {
			os.Exit(e.Code)
		} else {
			os.Exit(-1)
		}
	}
}
