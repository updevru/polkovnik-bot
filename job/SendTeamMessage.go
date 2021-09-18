package job

import (
	"polkovnik/adapter/notifyChannel"
	"polkovnik/domain"
)

func (p Processor) SendTeamMessage(team *domain.Team, task *domain.Task, story *domain.History, channel notifyChannel.Interface) error {
	settings := task.GetTaskSendTeamMessageSettingsDto()
	story.AddLine("Send message: " + settings.GetMessage())

	_, err := channel.SendTeamMessage(
		notifyChannel.Message{Text: settings.GetMessage()},
	)
	if err != nil {
		return err
	}

	story.AddLine("Message sent")

	return nil
}
