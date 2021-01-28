package notifyChannel

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

const TelegramChannelType = "telegram"

type TelegramChannel struct {
	token     string
	teamGroup int64
	api       *tgbotapi.BotAPI
}

func (t *TelegramChannel) SendTeamMessage(message Message) (bool, error) {
	msg := tgbotapi.NewMessage(t.teamGroup, message.Text)
	msg.ParseMode = "HTML"
	_, err := t.api.Send(msg)

	if err != nil {
		return false, err
	}

	log.Info("Send message to telegram", msg.Text)

	return true, nil
}

func NewTelegram(token string, teamGroup int64) (*TelegramChannel, error) {
	api, err := tgbotapi.NewBotAPI(token)

	return &TelegramChannel{token: token, teamGroup: teamGroup, api: api}, err
}
