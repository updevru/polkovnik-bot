package issueTracker

import (
	"fmt"
	jira "github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
	"net/url"
	"polkovnik/domain"
	"strings"
	"time"
)

const JiraTrackerTape = "jira"

type JiraTracker struct {
	url      url.URL
	username string
	password string
}

func NewJira(url url.URL, username string, password string) *JiraTracker {
	return &JiraTracker{
		url:      url,
		username: username,
		password: password,
	}
}

func (j JiraTracker) GetWorkLogByDate(date time.Time, projects []string) ([]WorkLogResponse, error) {

	dateFormatted := date.Format("2006/01/02")
	api, err := j.getApi()
	if err != nil {
		return nil, err
	}

	options := &jira.SearchOptions{Fields: []string{"*all"}}
	dql := fmt.Sprintf("project IN(%s) AND worklogDate = \"%s\"", strings.Join(projects, ","), dateFormatted)
	issues, _, err := api.Issue.Search(dql, options)
	if err != nil {
		return nil, err
	}

	var result []WorkLogResponse
	for _, issue := range issues {
		issueLogger := log.WithFields(log.Fields{"Ticket": issue.Key})
		issueLogger.Info(issue.Fields.Summary)

		for _, logItem := range issue.Fields.Worklog.Worklogs {
			logTime := time.Time(*logItem.Started)
			if dateFormatted != logTime.Format("2006/01/02") {
				continue
			}

			item := WorkLogResponse{
				Task:       Task{Id: issue.Key},
				User:       User{Login: logItem.Author.Name},
				LoggedTime: domain.NewTime(logItem.TimeSpentSeconds),
				Date:       logTime,
			}
			result = append(result, item)
			issueLogger.WithFields(log.Fields{
				"Author": logItem.Author.Name,
				"Time":   logItem.TimeSpent,
				"Date":   logTime.Format("2006/01/02 15:04:05"),
			}).Info("Logged time")
		}
	}

	return result, err
}

func (j JiraTracker) getApi() (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: j.username,
		Password: j.password,
	}

	return jira.NewClient(tp.Client(), j.url.String())
}
