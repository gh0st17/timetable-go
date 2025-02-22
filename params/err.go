package params

import "fmt"

var (
	ErrFewArguments     = fmt.Errorf("недостаточно аргументов")
	ErrMissingDep       = fmt.Errorf("номер института не указан")
	ErrMissingCourse    = fmt.Errorf("номер курса не указан")
	ErrWeekInput        = fmt.Errorf("номер недели должен быть от 1 до 18")
	ErrDepOutOfBound    = fmt.Errorf("номер института должен быть от 1 до 12 или 14")
	ErrCourseOutOfBound = fmt.Errorf("номер курса должен быть от 1 до 6")
)
