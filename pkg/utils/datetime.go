package utils

import "time"

func IsWeekend(date time.Time) bool {
	day := date.Weekday()
	return day == time.Saturday || day == time.Sunday
}

func CountWorkingDays(start, end time.Time) int {
	if start.After(end) {
		return 0
	}

	count := 0
	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		if !IsWeekend(d) {
			count++
		}
	}
	return count
}