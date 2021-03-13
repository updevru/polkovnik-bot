package domain

import (
	"github.com/google/uuid"
	"time"
)

type Team struct {
	Id           string
	Title        string
	Users        []*User
	Tasks        []*Task
	Channel      NotifyChannel
	Weekend      Weekend
	IssueTracker IssueTracker
	MinWorkLog   int
	DateCreated  time.Time
}

func NewTeam(title string) *Team {
	return &Team{
		Id:          uuid.NewString(),
		Title:       title,
		DateCreated: time.Now(),
	}
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
