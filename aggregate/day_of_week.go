package main

import "time"

func getDayOfWeekChar(date, format string) (dayOfWeek string, err error) {
	t, err := time.Parse(format, date)
	if err != nil {
		return "", err
	}

	switch t.Weekday() {
	case time.Monday:
		dayOfWeek = "月"
	case time.Tuesday:
		dayOfWeek = "火"
	case time.Wednesday:
		dayOfWeek = "水"
	case time.Thursday:
		dayOfWeek = "木"
	case time.Friday:
		dayOfWeek = "金"
	case time.Saturday:
		dayOfWeek = "土"
	default: // Sunday
		dayOfWeek = "日"
	}
	return dayOfWeek, nil
}
