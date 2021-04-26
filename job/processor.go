package job

import (
	log "github.com/sirupsen/logrus"
	"polkovnik/adapter/issueTracker"
	"polkovnik/adapter/notifyChannel"
	"polkovnik/app"
	"polkovnik/domain"
	"time"
)

type Processor struct {
	Tpl *app.TemplateEngine
}

func (p Processor) ProcessTeamTasks(team *domain.Team, date time.Time) error {
	if team.Weekend.IsWeekend(date) == true {
		log.Info("Team ", team.Title, " skip, is weekend")
		return nil
	}

	var err error
	tracker, err := issueTracker.New(team.IssueTracker)
	channel, err := notifyChannel.New(team.Channel)

	if err != nil {
		return err
	}

	for _, task := range team.Tasks {
		if !task.IsRun(date) {
			log.Info("Task skip ", task.Type, " last run ", task.LastRunTime)
			continue
		}

		log.Info("Task start ", task.Type)

		var err error
		switch task.Type {
		case domain.CheckTeamWorkLog:
			err = p.CheckTeamWorkLog(team, task, tracker, date, channel)
		case domain.SendTeamMessage:
			err = p.SendTeamMessage(team, task, channel)
		}

		if err != nil {
			return err
		}

		task.LastRunTime = time.Now().In(time.Local)
		log.Info("Task end ", task.Type)
	}

	return nil
}
