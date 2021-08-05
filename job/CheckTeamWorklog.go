package job

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"polkovnik/adapter/issueTracker"
	"polkovnik/adapter/notifyChannel"
	"polkovnik/domain"
	"time"
)

type teamMessageDataEntry struct {
	User       domain.User
	LoggedTime string
	NeedTime   int
	Points     int
	Difference int
}

type teamMessageData struct {
	List []teamMessageDataEntry
	Date string
}

func (p Processor) CheckTeamWorkLog(team *domain.Team, task *domain.Task, story *domain.History, tracker issueTracker.Interface, date time.Time, channel notifyChannel.Interface) error {
	var dateChek = date
	// Если в задании есть модификатор, то применяем его
	if len(task.DateModify) > 0 {
		duration, err := time.ParseDuration(task.DateModify)
		if err != nil {
			return err
		}
		dateChek = date.Add(duration)
	}

	story.AddLine("Fetch data from issue tracker by date - " + dateChek.Format("2006-01-02"))
	logged, err := issueTracker.CalculateTeamWorkLog(team, task, tracker, dateChek)
	if err != nil {
		return err
	}

	var data []teamMessageDataEntry
	for _, logEntry := range logged {
		if logEntry.User.Active == false {
			story.AddLine(fmt.Sprintf("User %s is not active", logEntry.User.Name))
			continue
		}

		if logEntry.User.Weekend.IsWeekend(dateChek) {
			story.AddLine(fmt.Sprintf("User %s on weekend", logEntry.User.Name))
			continue
		}

		story.AddLine(fmt.Sprintf("User %s logged time %s need %s", logEntry.User.Name, logEntry.Time.ToHumanFormat(), domain.NewTime(team.MinWorkLog).ToHumanFormat()))

		if team.MinWorkLog > logEntry.Time.Seconds() {
			data = append(
				data,
				teamMessageDataEntry{
					User:       logEntry.User,
					LoggedTime: logEntry.Time.ToHumanFormat(),
					NeedTime:   team.MinWorkLog,
					Points:     logEntry.User.Points,
					Difference: task.DeleteUserPoints,
				},
			)
		}
		log.Info(logEntry.User.Login, " logged time ", logEntry.Time.ToHumanFormat(), " need ", team.MinWorkLog)
	}

	if len(data) == 0 {
		story.AddLine("All users logged time!")
		return nil
	}

	message, err := channel.CreateMessageFromTemplate(
		"CheckTeamWorklog.html",
		teamMessageData{List: data, Date: dateChek.Format("02.01.2006")},
	)
	if err != nil {
		return err
	}

	story.AddLine("Send message: " + message.Text)
	_, err = channel.SendTeamMessage(message)
	if err != nil {
		return err
	}

	story.AddLine("Message sent")

	return nil
}
