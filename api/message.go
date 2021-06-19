package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"polkovnik/adapter/notifyChannel"
)

type messageRequest struct {
	Text string `json:"text"`
}

func (a apiHandler) MessageSend() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		team := a.store.GetTeam(vars["teamId"])
		if team == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Team #%s not found", vars["teamId"])})
			return
		}

		channel, err := notifyChannel.New(team.Channel)

		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
		}

		var request *messageRequest
		var body []byte

		body, err = ioutil.ReadAll(r.Body)
		if err == nil {
			err = json.Unmarshal(body, &request)
		}
		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
		}

		_, err = channel.SendTeamMessage(notifyChannel.Message{Text: request.Text})

		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
		}

		renderJson(w, http.StatusOK, ResponseSuccess{Result: "Ok"})
	})
}
