package domain

import (
	"errors"
	"github.com/google/uuid"
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
	Id               string
	Active           bool
	Schedule         *Schedule
	LastRunTime      time.Time
	Type             string   `validate:"required"`
	Projects         []string `validate:"required"`
	Message          string
	DateModify       string
	AddUserPoints    int
	DeleteUserPoints int
}

func validateTask(typeTask string, scheduleWeekDays []string, scheduleHour int, scheduleMinute int, projects []string, message string, dateModify string) error {
	if typeTask != CheckTeamWorkLog && typeTask != SendTeamMessage {
		return errors.New("invalid task")
	}

	if len(scheduleWeekDays) == 0 {
		return errors.New("schedule weekdays must be set")
	}

	if scheduleHour < 0 || scheduleHour > 24 {
		return errors.New("schedule hour must be more 0 and less 24")
	}

	if scheduleMinute < 0 || scheduleMinute > 60 {
		return errors.New("schedule minute must be more 0 and less 60")
	}

	if typeTask == SendTeamMessage {
		if len(message) < 3 {
			return errors.New("message must be set")
		}
	}

	if typeTask == CheckTeamWorkLog {
		if len(projects) == 0 {
			return errors.New("projects must be set")
		}

		for _, row := range projects {
			if len(row) < 2 {
				return errors.New("project name not be empty")
			}
		}

		if len(dateModify) > 0 {
			_, err := time.ParseDuration(dateModify)
			if err != nil {
				return errors.New("date modify error: " + err.Error())
			}
		}
	}

	return nil
}

func NewTask(typeTask string, scheduleWeekDays []string, scheduleHour int, scheduleMinute int, projects []string, message string, dateModify string) (*Task, error) {
	err := validateTask(typeTask, scheduleWeekDays, scheduleHour, scheduleMinute, projects, message, dateModify)
	if err != nil {
		return nil, err
	}

	return &Task{
		Id: uuid.NewString(),
		Schedule: &Schedule{
			WeekDays: scheduleWeekDays,
			Hour:     scheduleHour,
			Minute:   scheduleMinute,
		},
		Active:     true,
		Type:       typeTask,
		Projects:   projects,
		Message:    message,
		DateModify: dateModify,
	}, nil
}

func (t *Task) Edit(typeTask string, scheduleWeekDays []string, scheduleHour int, scheduleMinute int, projects []string, message string, dateModify string, active bool) error {
	err := validateTask(typeTask, scheduleWeekDays, scheduleHour, scheduleMinute, projects, message, dateModify)
	if err != nil {
		return err
	}

	t.Type = typeTask
	t.Projects = projects
	t.Message = message
	t.Schedule.WeekDays = scheduleWeekDays
	t.Schedule.Minute = scheduleMinute
	t.Schedule.Hour = scheduleHour
	t.DateModify = dateModify
	t.Active = active

	return nil
}

func (t Task) IsRun(date time.Time) bool {
	if t.Active == false {
		return false
	}

	next := t.Schedule.GetStartTime(date)
	if next == nil {
		return false
	}

	if next.Format("2006.01.02") == t.LastRunTime.Format("2006.01.02") {
		return false
	}

	return date.After(*next)
}
