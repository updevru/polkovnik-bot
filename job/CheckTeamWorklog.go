package job

import (
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

func (p Processor) CheckTeamWorkLog(team *domain.Team, task *domain.Task, tracker issueTracker.Interface, date time.Time, channel notifyChannel.Interface) error {
	var dateChek = date
	// Если в задании есть модификатор, то применяем его
	if len(task.DateModify) > 0 {
		duration, err := time.ParseDuration(task.DateModify)
		if err != nil {
			return err
		}
		dateChek = date.Add(duration)
	}

	logged, err := issueTracker.CalculateTeamWorkLog(team, task, tracker, dateChek)
	if err != nil {
		return err
	}

	var data []teamMessageDataEntry
	for _, logEntry := range logged {
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
		return nil
	}

	message, err := p.Tpl.RenderString(
		"telegram/CheckTeamWorklog.html",
		teamMessageData{List: data, Date: dateChek.Format("02.01.2006")},
	)
	if err != nil {
		return err
	}

	_, err = channel.SendTeamMessage(
		notifyChannel.Message{Text: message},
	)
	if err != nil {
		return err
	}

	return nil
}
