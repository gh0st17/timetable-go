package database

import "fmt"

var (
	ErrOpenDB     = fmt.Errorf("ошибка открытия базы данных")
	ErrCloseDB    = fmt.Errorf("ошибка закрытия базы данных")
	ErrInsetGroup = fmt.Errorf("ошибка при добавлении записи в базу данных")
	ErrQueryGroup = fmt.Errorf("ошибка при запросе групп в базе данных")
	ErrReadGroups = fmt.Errorf("ошибка чтения базы данных")
	ErrDelete     = fmt.Errorf("ошибка при запросе удаления групп в базе данных")
)
