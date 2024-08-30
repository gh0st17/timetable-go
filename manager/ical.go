package manager

import "time"

func calcWeek() uint8 {
	today := time.Now()
	_, week := today.ISOWeek()

	if today.Month() >= 8 && today.Day() >= 1 {
		week -= 34
	} else {
		week -= 6
	}

	if week < 1 {
		week = 1
	} else if week > 18 {
		week = 18
	}

	return uint8(week)
}
