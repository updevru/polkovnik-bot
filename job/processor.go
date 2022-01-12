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

type messageQueue struct {
	task *domain.Task
	team *domain.Team
	date time.Time
}

type Processor struct {
	Tpl     *app.TemplateEngine
	config  *domain.Config
	history *repository.HistoryRepository
	queue   chan *messageQueue
	ticker  *time.Ticker
	lock    map[string]string
}

func NewProcessor(tpl *app.TemplateEngine, config *domain.Config, history *repository.HistoryRepository) *Processor {
	return &Processor{
		Tpl:     tpl,
		config:  config,
		history: history,
		lock:    make(map[string]string, 10),
		queue:   make(chan *messageQueue, 20),
	}
}

func (p *Processor) StartScheduler() {
	p.ticker = time.NewTicker(time.Minute)

	for tick := range p.ticker.C {
		now := tick.In(time.Local)
		for _, team := range p.config.Teams {
			log.Info("Run team ", team.Title)

			if team.Weekend.IsWeekend(now) == true {
				log.Info("Team ", team.Title, " skip, is weekend")
				continue
			}

			for _, task := range team.Tasks {
				if !task.IsRun(now) {
					log.Info("Task skip ", task.Type, " last run ", task.LastRunTime, " active ", task.Active)
					continue
				}

				//Возможно задача уже выполняется или ждет своей очереди
				if p.isLockTask(task) {
					continue
				}

				p.lockTask(task)
				p.ScheduleTask(team, task, now)
			}
		}
	}
}

func (p *Processor) Stop() {
	p.ticker.Stop()
	close(p.queue)
}

func (p *Processor) ScheduleTask(team *domain.Team, task *domain.Task, date time.Time) {
	log.Info("Task schedule ", task.Type)
	message := &messageQueue{
		task: task,
		team: team,
		date: date,
	}
	p.queue <- message
	log.Info("Task queue len = ", len(p.queue))
}

func (p *Processor) StartWorker() {
	for message := range p.queue {
		task := message.task
		team := message.team

		log.Info("Task start ", task.Type)
		story := domain.NewHistory(task.Id)

		var err error
		tracker, err := issueTracker.New(team.IssueTracker)
		channel, err := notifyChannel.New(team.Channel, p.Tpl)

		if err != nil {
			story.SetError()
			story.AddLine("Error: " + err.Error())
			log.Error("Error: ", err.Error())
			p.history.New(story)
			p.unlockTask(task)
			continue
		}

		switch task.Type {
		case domain.CheckTeamWorkLog:
			err = p.CheckTeamWorkLog(team, task, story, tracker, message.date, channel)
		case domain.CheckTeamWorkLogByPeriod:
			err = p.CheckTeamWorkLogByPeriod(team, task, story, tracker, message.date, channel)
		case domain.SendTeamMessage:
			err = p.SendTeamMessage(team, task, story, channel)
		case domain.CheckUserWeekend:
			err = p.CheckUserWeekend(team, task, story, channel, message.date)
		}

		if err != nil {
			story.SetError()
			story.AddLine("Error: " + err.Error())
			p.history.New(story)
			p.unlockTask(task)
			continue
		}

		task.LastRunTime = time.Now().In(time.Local)
		log.Info("Task end ", task.Type)
		story.SetSuccess()
		story.AddLine("Task completed")

		p.history.New(story)
		p.unlockTask(task)
	}
}

func (p *Processor) lockTask(task *domain.Task) {
	p.lock[task.Id] = task.Id
}

func (p *Processor) isLockTask(task *domain.Task) bool {
	_, found := p.lock[task.Id]
	return found
}

func (p *Processor) unlockTask(task *domain.Task) {
	if p.isLockTask(task) {
		delete(p.lock, task.Id)
	}
}
