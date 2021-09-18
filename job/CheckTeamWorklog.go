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
}

type teamMessageData struct {
	List []teamMessageDataEntry
	Date string
}

type teamMessageDataByDate struct {
	Periods []teamMessageData
}

func (p Processor) CheckTeamWorkLog(team *domain.Team, task *domain.Task, story *domain.History, tracker issueTracker.Interface, date time.Time, channel notifyChannel.Interface) error {
	settings := task.GetTaskCheckTeamWorkLogSettingsDto()
	var dateChek = date
	// Если в задании есть модификатор, то применяем его
	if len(settings.GetDateModifyDuration()) > 0 {
		duration, err := time.ParseDuration(settings.GetDateModifyDuration())
		if err != nil {
			return err
		}
		dateChek = date.Add(duration)
	}

	data, err := getLoggedByDate(team, settings.GetProjectsList(), story, tracker, dateChek)
	if err != nil {
		return err
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

	story.AddLine("Send message:")
	story.AddLine(message.Text)
	_, err = channel.SendTeamMessage(message)
	if err != nil {
		return err
	}

	story.AddLine("Message sent")

	return nil
}

func (p Processor) CheckTeamWorkLogByPeriod(team *domain.Team, task *domain.Task, story *domain.History, tracker issueTracker.Interface, date time.Time, channel notifyChannel.Interface) error {
	settings := task.GetTaskCheckTeamWorkLogByPeriodSettingsDto()

	var dateStart = time.Now().In(time.Local)
	duration, _ := time.ParseDuration(settings.GetStartDuration())
	dateStart = dateStart.Add(duration)

	var dateEnd = time.Now().In(time.Local)
	duration, _ = time.ParseDuration(settings.GetEndDuration())
	dateEnd = dateEnd.Add(duration)

	dateChek := dateStart
	var result = make([]teamMessageData, 0)
	for {
		if dateEnd.Before(dateChek) {
			break
		}

		if team.Weekend.IsWeekend(dateChek) == true {
			story.AddLine(fmt.Sprintf("On date %s team is weekend, skip", dateChek.String()))
			dateChek = dateChek.AddDate(0, 0, 1)
			continue
		}

		data, err := getLoggedByDate(team, settings.GetProjectsList(), story, tracker, dateChek)
		if err != nil {
			return err
		}

		if len(data) == 0 {
			story.AddLine(fmt.Sprintf("By date %s all users logged time!", dateChek.String()))
			dateChek = dateChek.AddDate(0, 0, 1)
			continue
		}

		result = append(result, teamMessageData{List: data, Date: dateChek.Format("02.01.2006")})
		dateChek = dateChek.AddDate(0, 0, 1)
	}

	if len(result) == 0 {
		story.AddLine("All users logged time!")
		return nil
	}

	message, err := channel.CreateMessageFromTemplate(
		"CheckTeamWorklogByPeriod.html",
		teamMessageDataByDate{Periods: result},
	)
	if err != nil {
		return err
	}

	story.AddLine("Send message:")
	story.AddLine(message.Text)
	_, err = channel.SendTeamMessage(message)
	if err != nil {
		return err
	}

	story.AddLine("Message sent")

	return nil
}

func getLoggedByDate(team *domain.Team, projects []string, story *domain.History, tracker issueTracker.Interface, dateChek time.Time) ([]teamMessageDataEntry, error) {
	story.AddLine("Fetch data from issue tracker by date - " + dateChek.Format("2006-01-02"))
	logged, err := issueTracker.CalculateTeamWorkLog(team, projects, tracker, dateChek)
	if err != nil {
		return nil, err
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
				},
			)
		}
		log.Info(logEntry.User.Login, " logged time ", logEntry.Time.ToHumanFormat(), " need ", team.MinWorkLog)
	}
	return data, err
}
