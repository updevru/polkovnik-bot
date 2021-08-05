package notifyChannel

import (
	"errors"
	"polkovnik/app"
	"polkovnik/domain"
	"strconv"
)

type Message struct {
	Text string
	Data interface{}
}

type Interface interface {
	SendTeamMessage(message Message) (bool, error)
	CreateMessageFromTemplate(template string, data interface{}) (Message, error)
}

func New(channel *domain.NotifyChannel, tpl *app.TemplateEngine) (Interface, error) {
	if len(channel.Type) == 0 {
		return nil, errors.New("channel type not defined")
	}

	var result Interface
	var err error

	switch channel.Type {
	case TelegramChannelType:
		var id int64
		id, err = strconv.ParseInt(channel.ChannelId, 10, 64)

		if err == nil {
			result, err = NewTelegram(channel.Settings["token"], id, tpl)
		}

		break
	case WebexChannelType:
		result, err = NewWebex(channel.Settings["token"], channel.ChannelId, tpl)
		break
	default:
		return nil, errors.New("channel type not found")
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
