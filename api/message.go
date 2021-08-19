package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"polkovnik/adapter/notifyChannel"
)

// Сообщение для отправки команде.
// swagger:parameters MessageSend
type messageRequestWrapper struct {
	// in:body
	Body messageRequest
}

// Сообщение
// swagger:model Message
type messageRequest struct {
	//Текст сообщения
	Text string `json:"text"`
}

func (m messageRequest) isValid() error {

	if len(m.Text) < 1 {
		return errors.New("message is empty")
	}

	return nil
}

// swagger:route POST /team/{teamId}/sendMessage Messages MessageSend
//
// Отправка сообщения команде.
//
// Responses:
//        200: ResponseSuccess
//        400: ResponseError
//        404: ResponseError
func (a apiHandler) MessageSend() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		team := a.store.GetTeam(vars["teamId"])
		if team == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Team #%s not found", vars["teamId"])})
			return
		}

		channel, err := notifyChannel.New(team.Channel, a.processor.Tpl)

		var request *messageRequest
		var body []byte

		if err == nil {
			body, err = ioutil.ReadAll(r.Body)
		}

		if err == nil {
			err = json.Unmarshal(body, &request)
		}

		if err == nil {
			err = request.isValid()
		}

		if err == nil {
			_, err = channel.SendTeamMessage(notifyChannel.Message{Text: request.Text})
		}

		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		renderJson(w, http.StatusOK, ResponseSuccess{Result: "Ok"})
	})
}
