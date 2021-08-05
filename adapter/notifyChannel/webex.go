package notifyChannel

import (
	webex "github.com/jbogarin/go-cisco-webex-teams/sdk"
	log "github.com/sirupsen/logrus"
	"polkovnik/app"
)

const WebexChannelType = "webex"

type WebexChannel struct {
	token     string
	teamGroup string
	api       *webex.Client
	tpl       *app.TemplateEngine
}

func (w *WebexChannel) SendTeamMessage(message Message) (bool, error) {
	msg := &webex.MessageCreateRequest{
		RoomID:   w.teamGroup,
		Markdown: message.Text,
	}
	_, _, err := w.api.Messages.CreateMessage(msg)

	if err != nil {
		return false, err
	}

	log.Info("Send message to webex", msg.Markdown)

	return true, nil
}

func (w *WebexChannel) CreateMessageFromTemplate(template string, data interface{}) (Message, error) {
	body, err := w.tpl.RenderString("webex/"+template, data)

	return Message{
		Text: body,
	}, err
}

func NewWebex(token string, teamGroup string, tpl *app.TemplateEngine) (*WebexChannel, error) {

	api := webex.NewClient()
	api.SetAuthToken(token)

	return &WebexChannel{
		token:     token,
		teamGroup: teamGroup,
		api:       api,
		tpl:       tpl,
	}, nil
}
