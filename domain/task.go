package domain

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

const CheckTeamWorkLog = "check_work_log"
const CheckTeamWorkLogByPeriod = "check_work_log_by_period"
const SendTeamMessage = "send_team_message"
const CheckUserWeekend = "check_user_weekend"

var taskTypes = []string{CheckTeamWorkLog, CheckTeamWorkLogByPeriod, SendTeamMessage, CheckUserWeekend}

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

type TaskSettings map[string]string

type TaskSettingsInterface interface {
	GetSettings() TaskSettings
	Validate() error
}

type Task struct {
	Id               string
	Active           bool
	Schedule         *Schedule
	LastRunTime      time.Time
	Type             string
	Projects         []string
	Message          string
	DateModify       string
	AddUserPoints    int
	DeleteUserPoints int
	TaskSettings
}

func (t *Task) setSettings(settings TaskSettingsInterface) error {
	err := settings.Validate()
	if err != nil {
		return err
	}

	t.TaskSettings = settings.GetSettings()

	return nil
}

func (t Task) GetTaskCheckTeamWorkLogByPeriodSettingsDto() TaskCheckTeamWorkLogByPeriodSettingsDto {
	return TaskCheckTeamWorkLogByPeriodSettingsDto{TaskSettingsDto{TaskSettings: t.TaskSettings}}
}

func (t Task) GetTaskCheckUserWeekendSettingsDto() TaskCheckUserWeekendSettingsDto {
	return TaskCheckUserWeekendSettingsDto{TaskSettingsDto{TaskSettings: t.TaskSettings}}
}

func (t Task) GetTaskCheckTeamWorkLogSettingsDto() TaskCheckTeamWorkLogSettingsDto {
	return TaskCheckTeamWorkLogSettingsDto{TaskSettingsDto{TaskSettings: t.TaskSettings}}
}

func (t Task) GetTaskSendTeamMessageSettingsDto() TaskSendTeamMessageSettingsDto {
	return TaskSendTeamMessageSettingsDto{TaskSettingsDto{TaskSettings: t.TaskSettings}}
}

func createTaskSettingsDto(typeTask string, settings map[string]string) TaskSettingsInterface {
	dto := TaskSettingsDto{TaskSettings: settings}
	switch typeTask {
	case CheckTeamWorkLog:
		return TaskCheckTeamWorkLogSettingsDto{dto}
	case CheckTeamWorkLogByPeriod:
		return TaskCheckTeamWorkLogByPeriodSettingsDto{dto}
	case SendTeamMessage:
		return TaskSendTeamMessageSettingsDto{dto}
	case CheckUserWeekend:
		return TaskCheckUserWeekendSettingsDto{dto}
	}

	return nil
}

type TaskSettingsDto struct {
	TaskSettings
}

func (t TaskSettingsDto) GetSettings() TaskSettings {
	return t.TaskSettings
}

type TaskCheckUserWeekendSettingsDto struct {
	TaskSettingsDto
}

func (t TaskCheckUserWeekendSettingsDto) GetDateModifyDuration() string {
	return t.TaskSettings["date_modify"]
}

func (t TaskCheckUserWeekendSettingsDto) Validate() error {

	if len(t.GetDateModifyDuration()) == 0 {
		return errors.New("date_modify must be set")
	}

	_, err := time.ParseDuration(t.GetDateModifyDuration())
	if err != nil {
		return errors.New("date modify error: " + err.Error())
	}

	return nil
}

type TaskCheckTeamWorkLogSettingsDto struct {
	TaskSettingsDto
}

func (t TaskCheckTeamWorkLogSettingsDto) GetProjects() string {
	return t.TaskSettings["projects"]
}

func (t TaskCheckTeamWorkLogSettingsDto) GetProjectsList() []string {
	return strings.Split(t.TaskSettings["projects"], ",")
}

func (t TaskCheckTeamWorkLogSettingsDto) GetDateModifyDuration() string {
	return t.TaskSettings["date_modify"]
}

func (t TaskCheckTeamWorkLogSettingsDto) Validate() error {

	if len(t.GetProjects()) == 0 {
		return errors.New("projects must be set")
	}

	if len(t.GetDateModifyDuration()) > 0 {
		_, err := time.ParseDuration(t.GetDateModifyDuration())
		if err != nil {
			return errors.New("date modify error: " + err.Error())
		}
	}

	return nil
}

type TaskCheckTeamWorkLogByPeriodSettingsDto struct {
	TaskSettingsDto
}

func (t TaskCheckTeamWorkLogByPeriodSettingsDto) GetStartDuration() string {
	return t.TaskSettings["start_duration"]
}

func (t TaskCheckTeamWorkLogByPeriodSettingsDto) GetEndDuration() string {
	return t.TaskSettings["end_duration"]
}

func (t TaskCheckTeamWorkLogByPeriodSettingsDto) GetProjects() string {
	return t.TaskSettings["projects"]
}

func (t TaskCheckTeamWorkLogByPeriodSettingsDto) GetProjectsList() []string {
	return strings.Split(t.TaskSettings["projects"], ",")
}

func (t TaskCheckTeamWorkLogByPeriodSettingsDto) Validate() error {
	if len(t.GetProjects()) == 0 {
		return errors.New("projects must be set")
	}

	var err error
	var start time.Duration
	if len(t.GetStartDuration()) == 0 {
		return errors.New("start_duration must be set")
	} else {
		start, err = time.ParseDuration(t.GetStartDuration())
		if err != nil {
			return errors.New("start_duration modify error: " + err.Error())
		}

		if start.Seconds() > 0 {
			return errors.New("start_duration must be in the past")
		}
	}

	var end time.Duration
	if len(t.GetEndDuration()) == 0 {
		return errors.New("end_duration must be set")
	} else {
		end, err = time.ParseDuration(t.GetEndDuration())
		if err != nil {
			return errors.New("end_duration modify error: " + err.Error())
		}

		if end.Seconds() > 0 {
			return errors.New("end_duration must be in the past")
		}
	}

	if start.Seconds() > end.Seconds() {
		return errors.New("start_duration must be less than end_duration")
	}

	return nil
}

type TaskSendTeamMessageSettingsDto struct {
	TaskSettingsDto
}

func (t TaskSendTeamMessageSettingsDto) GetMessage() string {
	return t.TaskSettings["message"]
}

func (t TaskSendTeamMessageSettingsDto) Validate() error {
	if len(t.GetMessage()) < 3 {
		return errors.New("message must be set")
	}

	return nil
}

func validateTask(typeTask string, scheduleWeekDays []string, scheduleHour int, scheduleMinute int, settings TaskSettingsInterface) error {

	var typeExist = false
	for _, v := range taskTypes {
		if v == typeTask {
			typeExist = true
		}
	}

	if typeExist == false {
		return errors.New("wrong task type")
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

	err := settings.Validate()
	if err != nil {
		return err
	}

	return nil
}

func NewTask(typeTask string, scheduleWeekDays []string, scheduleHour int, scheduleMinute int, settings map[string]string) (*Task, error) {
	settingsDto := createTaskSettingsDto(typeTask, settings)
	if settingsDto == nil {
		return nil, errors.New("task type not found")
	}

	err := validateTask(typeTask, scheduleWeekDays, scheduleHour, scheduleMinute, settingsDto)
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
		Active:       true,
		Type:         typeTask,
		TaskSettings: settingsDto.GetSettings(),
	}, nil
}

func (t *Task) Edit(typeTask string, scheduleWeekDays []string, scheduleHour int, scheduleMinute int, active bool, settings map[string]string) error {
	settingsDto := createTaskSettingsDto(typeTask, settings)
	if settingsDto == nil {
		return errors.New("task type not found")
	}

	err := validateTask(typeTask, scheduleWeekDays, scheduleHour, scheduleMinute, settingsDto)
	if err != nil {
		return err
	}

	t.Type = typeTask
	t.Schedule.WeekDays = scheduleWeekDays
	t.Schedule.Minute = scheduleMinute
	t.Schedule.Hour = scheduleHour
	t.Active = active
	t.setSettings(settingsDto)

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
