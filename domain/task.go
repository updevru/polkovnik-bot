package domain

import (
	"time"
)

const CheckTeamWorkLog = "check_work_log"
const SendTeamMessage = "send_team_message"

type Schedule struct {
	WeekDays []string
	Hour     int
	Minute   int
}

func (s Schedule) isThisDay(date time.Time) bool {
	day := date.Weekday().String()
	for _, row := range s.WeekDays {
		if row == day {
			return true
		}
	}

	return false
}

func (s Schedule) GetStartTime(date time.Time) *time.Time {
	if !s.isThisDay(date) {
		return nil
	}

	next := time.Date(date.Year(), date.Month(), date.Day(), s.Hour, s.Minute, 0, 0, time.Local)

	return &next
}

type Task struct {
	Schedule         Schedule
	LastRunTime      time.Time
	Type             string
	Projects         []string
	Message          string
	DateModify       string
	AddUserPoints    int
	DeleteUserPoints int
}

func (t Task) IsRun(date time.Time) bool {
	next := t.Schedule.GetStartTime(date)
	if next == nil {
		return false
	}

	if next.Format("2006.01.02") == t.LastRunTime.Format("2006.01.02") {
		return false
	}

	return date.After(*next)
}
