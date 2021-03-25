package notifyChannel

import (
	"polkovnik/domain"
	"strconv"
)

type Message struct {
	Text string
}

type Interface interface {
	SendTeamMessage(message Message) (bool, error)
}

func New(channel *domain.NotifyChannel) Interface {
	var result Interface
	switch channel.Type {
	case TelegramChannelType:
		var err error
		id, _ := strconv.ParseInt(channel.ChannelId, 10, 64)

		result, err = NewTelegram(channel.Settings["token"], id)

		if err != nil {
			panic(err.Error())
		}

		break
	default:
		panic("Channel type not found")
	}

	return result
}
