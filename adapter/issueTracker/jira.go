package issueTracker

import (
	"fmt"
	jira "github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"teamBot/domain"
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
	api, _ := j.getApi()
	options := &jira.SearchOptions{Fields: []string{"*all"}}
	dql := fmt.Sprintf("project IN(%s) AND worklogDate = \"%s\"", strings.Join(projects, ","), dateFormatted)
	issues, _, err := api.Issue.Search(dql, options)

	var result []WorkLogResponse
	for _, issue := range issues {
		log.Info("Ticket", issue.Key, issue.Fields.Summary)

		for _, log := range issue.Fields.Worklog.Worklogs {
			logTime := time.Time(*log.Created)
			if dateFormatted != logTime.Format("2006/01/02") {
				continue
			}

			item := WorkLogResponse{
				Task:       Task{Id: issue.Key},
				User:       User{Login: log.Author.Name},
				LoggedTime: domain.NewTime(log.TimeSpentSeconds),
				Date:       logTime,
			}
			result = append(result, item)
			//fmt.Println("Log Author", log.Author.Name, log.Author.EmailAddress ,"Time", log.TimeSpent, log.TimeSpentSeconds, "Date", logTime.Format("2006/01/02"))
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
