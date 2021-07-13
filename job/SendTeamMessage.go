package job

import (
	"errors"
	"polkovnik/adapter/notifyChannel"
	"polkovnik/domain"
)

func (p Processor) SendTeamMessage(team *domain.Team, task *domain.Task, story *domain.History, channel notifyChannel.Interface) error {
	if len(task.Message) <= 0 {
		return errors.New("message is empty")
	}

	story.AddLine("Send message: " + task.Message)

	_, err := channel.SendTeamMessage(
		notifyChannel.Message{Text: task.Message},
	)
	if err != nil {
		return err
	}

	story.AddLine("Message sent")

	return nil
}
