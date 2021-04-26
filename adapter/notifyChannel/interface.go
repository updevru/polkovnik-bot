package notifyChannel

import (
	"errors"
	"polkovnik/domain"
	"strconv"
)

type Message struct {
	Text string
}

type Interface interface {
	SendTeamMessage(message Message) (bool, error)
}

func New(channel *domain.NotifyChannel) (Interface, error) {
	var result Interface
	switch channel.Type {
	case TelegramChannelType:
		var err error
		id, err := strconv.ParseInt(channel.ChannelId, 10, 64)

		if err == nil {
			result, err = NewTelegram(channel.Settings["token"], id)
		}

		if err != nil {
			return nil, err
		}

		break
	default:
		return nil, errors.New("channel type not found")
	}

	return result, nil
}
