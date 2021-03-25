package issueTracker

import (
	"net/url"
	"polkovnik/domain"
	"time"
)

type WorkLogResponse struct {
	Task
	User
	LoggedTime domain.Time
	Date       time.Time
}

type Task struct {
	Id string
}

type User struct {
	Login string
}

type Interface interface {
	GetWorkLogByDate(time time.Time, projects []string) ([]WorkLogResponse, error)
}

func New(tracker *domain.IssueTracker) Interface {
	var IssueTracker Interface
	switch tracker.Type {
	case JiraTrackerTape:
		host, _ := url.Parse(tracker.Settings["url"])
		IssueTracker = NewJira(
			*host,
			tracker.Settings["username"],
			tracker.Settings["password"],
		)
	default:
		panic("Tracker type not found")
	}

	return IssueTracker
}

type CalculateTeamWorkLogResponse struct {
	User domain.User
	Time domain.Time
}

func CalculateTeamWorkLog(team *domain.Team, task *domain.Task, tracker Interface, date time.Time) ([]CalculateTeamWorkLogResponse, error) {
	var result = make(map[string]int, 20)
	logs, err := tracker.GetWorkLogByDate(date, task.Projects)
	if err != nil {
		return nil, err
	}

	for _, log := range logs {
		if _, exist := result[log.User.Login]; exist {
			result[log.User.Login] += log.LoggedTime.Seconds()
		} else {
			result[log.User.Login] = log.LoggedTime.Seconds()
		}
	}

	var response []CalculateTeamWorkLogResponse
	for _, user := range team.Users {
		var loggedTime int = 0
		if _, exist := result[user.Login]; exist {
			loggedTime = result[user.Login]
		}

		item := CalculateTeamWorkLogResponse{
			User: *user,
			Time: domain.NewTime(loggedTime),
		}
		response = append(response, item)
	}

	return response, nil
}
