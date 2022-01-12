package receiver

import (
	"polkovnik/adapter/notifyChannel"
	"polkovnik/domain"
)

func (p Processor) sendTeamMessage(team *domain.Team, receiver *domain.Receiver, dto TemplateDto) error {

	text, err := renderData(receiver.Settings["message"], dto)

	if err != nil {
		return err
	}

	channel, err := notifyChannel.New(team.Channel, p.Tpl)
	if err != nil {
		return err
	}

	_, err = channel.SendTeamMessage(notifyChannel.Message{Text: text})

	return err
}
