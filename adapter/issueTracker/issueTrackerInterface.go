package issueTracker

import (
	"errors"
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

func New(tracker *domain.IssueTracker) (Interface, error) {
	if len(tracker.Type) == 0 {
		return nil, errors.New("tracker type not defined")
	}

	var IssueTracker Interface
	switch tracker.Type {
	case JiraTrackerTape:
		host, err := url.Parse(tracker.Settings["url"])
		if err != nil {
			return nil, err
		}
		IssueTracker = NewJira(
			*host,
			tracker.Settings["username"],
			tracker.Settings["password"],
		)
	default:
		return nil, errors.New("tracker type not found")
	}

	return IssueTracker, nil
}

func GetPublicSettings(tracker *domain.IssueTracker) (map[string]string, error) {
	result := make(map[string]string)

	switch tracker.Type {
	case JiraTrackerTape:
		if val, ok := tracker.Settings["url"]; ok {
			result["url"] = val
		} else {
			result["url"] = ""
		}

		if val, ok := tracker.Settings["username"]; ok {
			result["username"] = val
		} else {
			result["username"] = ""
		}
	default:
		return nil, errors.New("tracker type not found")
	}

	return result, nil
}

type CalculateTeamWorkLogResponse struct {
	User domain.User
	Time domain.Time
}

func CalculateTeamWorkLog(team *domain.Team, projects []string, tracker Interface, date time.Time) ([]CalculateTeamWorkLogResponse, error) {
	var result = make(map[string]int, 20)
	logs, err := tracker.GetWorkLogByDate(date, projects)
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
