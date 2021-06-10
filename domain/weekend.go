package domain

import "time"

type Weekend struct {
	WeekDays  []string
	Intervals []WeekendInterval
}

func (w Weekend) IsWeekend(date time.Time) bool {

	for _, day := range w.WeekDays {
		if date.Weekday().String() == day {
			return true
		}
	}

	for _, interval := range w.Intervals {
		if interval.IsWeekend(date) {
			return true
		}
	}

	return false
}

type WeekendInterval struct {
	Start time.Time
	End   time.Time
}

func (i WeekendInterval) IsWeekend(date time.Time) bool {
	if date.Equal(i.Start) || date.Equal(i.End) {
		return true
	}

	if date.After(i.Start) && date.Before(i.End) {
		return true
	}

	if i.Start.Equal(i.End) {
		return date.Format("2006-01-02") == i.Start.Format("2006-01-02")
	}

	if date.Format("2006-01-02") == i.End.Format("2006-01-02") {
		return true
	}

	return false
}
