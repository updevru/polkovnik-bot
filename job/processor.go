package job

import (
	log "github.com/sirupsen/logrus"
	"polkovnik/adapter/issueTracker"
	"polkovnik/adapter/notifyChannel"
	"polkovnik/app"
	"polkovnik/domain"
	"polkovnik/repository"
	"time"
)

type Processor struct {
	Tpl *app.TemplateEngine
}

func (p Processor) ProcessTeamTasks(team *domain.Team, history *repository.HistoryRepository, date time.Time) error {
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
			log.Info("Task skip ", task.Type, " last run ", task.LastRunTime, " active ", task.Active)
			continue
		}

		log.Info("Task start ", task.Type)
		story := domain.NewHistory(task.Id)

		var err error
		switch task.Type {
		case domain.CheckTeamWorkLog:
			err = p.CheckTeamWorkLog(team, task, story, tracker, date, channel)
		case domain.SendTeamMessage:
			err = p.SendTeamMessage(team, task, story, channel)
		}

		if err != nil {
			story.SetError()
			story.AddLine("Error: " + err.Error())
			history.New(story)
			return err
		}

		task.LastRunTime = time.Now().In(time.Local)
		log.Info("Task end ", task.Type)
		story.SetSuccess()
		story.AddLine("Task completed")

		return history.New(story)
	}

	return nil
}
