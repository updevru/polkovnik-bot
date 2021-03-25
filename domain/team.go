package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Team struct {
	Id           string
	Title        string
	Users        []*User
	Tasks        []*Task
	Channel      *NotifyChannel
	Weekend      Weekend
	IssueTracker *IssueTracker
	MinWorkLog   int
	DateCreated  time.Time
}

func validateTeam(title string) error {

	if len(title) == 0 {
		return errors.New("title must be set")
	}

	return nil
}

func NewTeam(title string) (*Team, error) {
	err := validateTeam(title)
	if err != nil {
		return nil, err
	}

	return &Team{
		Id:          uuid.NewString(),
		Title:       title,
		DateCreated: time.Now(),
	}, nil
}

func validateTeamSettings(title string, notifyChannelType string, notifyChannelChannelId string, notifyChannelSettings map[string]string, issueTrackerType string, issueTrackerSettings map[string]string, minWorkLog int) error {

	if len(title) == 0 {
		return errors.New("title must be set")
	}

	if notifyChannelType != "telegram" {
		return errors.New("wrong channel type")
	}

	if len(notifyChannelChannelId) == 0 {
		return errors.New("notify channel ID must be set")
	}

	for name, value := range notifyChannelSettings {
		if len(value) == 0 {
			return errors.New(name + " must be set")
		}
	}

	if issueTrackerType != "jira" {
		return errors.New("wrong issue tracker type")
	}

	for name, value := range issueTrackerSettings {
		if len(value) == 0 {
			return errors.New(name + " must be set")
		}
	}

	if minWorkLog <= 0 {
		return errors.New("min work log must be more 0")
	}

	return nil
}

func (t *Team) EditSettings(title string, notifyChannelType string, notifyChannelChannelId string, notifyChannelSettings map[string]string, issueTrackerType string, issueTrackerSettings map[string]string, minWorkLog int) error {
	err := validateTeamSettings(title, notifyChannelType, notifyChannelChannelId, notifyChannelSettings, issueTrackerType, issueTrackerSettings, minWorkLog)
	if err != nil {
		return err
	}

	t.Title = title
	if t.Channel == nil {
		t.Channel = &NotifyChannel{}
	}
	t.Channel.Type = notifyChannelType
	t.Channel.ChannelId = notifyChannelChannelId
	t.Channel.Settings = notifyChannelSettings

	if t.IssueTracker == nil {
		t.IssueTracker = &IssueTracker{}
	}
	t.IssueTracker.Type = issueTrackerType
	t.IssueTracker.Settings = issueTrackerSettings
	t.MinWorkLog = minWorkLog

	return nil
}

func (t *Team) AddUserPoint(user User, point int) bool {
	for _, item := range t.Users {
		if item.Login == user.Login {
			item.AddPoint(point)
			return true
		}
	}

	return false
}

func (t *Team) DeleteUserPoint(user User, point int) bool {
	for _, item := range t.Users {
		if item.Login == user.Login {
			item.DeletePoint(point)
			return true
		}
	}

	return false
}
