package api

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"polkovnik/adapter/notifyChannel"
	"polkovnik/domain"
	"text/template"
)

type receiverTemplateDto struct {
	Method string
	Body   interface{}
	Params url.Values
	Header http.Header
}

// swagger:parameters ReceiversGet ReceiversEdit ReceiversDelete Receive
type receiverId struct {
	// ID Приемника
	// in: path
	// required: true
	ReceiverId string `json:"receiverId"`
}

// Приемник
// swagger:model Receiver
type receiverResponseItem struct {
	Id       string            `json:"id"`
	Active   bool              `json:"active"`
	Type     string            `json:"type"`
	Format   string            `json:"format"`
	Settings map[string]string `json:"settings"`
}

// swagger:response ReceiversGet
type receiverResponseItemWrapper struct {
	// in: body
	Body receiverResponseItem `json:"body"`
}

// swagger:parameters ReceiversAdd ReceiversEdit
type receiverItemRequestWrapper struct {
	// in: body
	Body receiverResponseItem `json:"body"`
}

//Список приемников
// swagger:model Receivers
type receiverResponseList struct {
	Result []receiverResponseItem `json:"result"`
}

// swagger:response ReceiversList
type receiverResponseListWrapper struct {
	// in: body
	Body receiverResponseList `json:"body"`
}

// Данные для приема
// swagger:parameters Receive
type receiverRequestWrapper struct {
	// in:body
	Body interface{}
}

// swagger:route GET /team/{teamId}/receive/{receiverId} Receivers Receive
//
// Прием данных приемником (поддерживает любые методы).
//
// Responses:
//        200: ResponseSuccess
//        400: ResponseError
//        404: ResponseError
func (a apiHandler) Receive() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		dump, _ := httputil.DumpRequest(r, true)
		fmt.Printf("Request: %q", dump)

		team := a.store.GetTeam(vars["teamId"])
		if team == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Team #%s not found", vars["teamId"])})
			return
		}

		var err error
		var post *interface{}
		var body []byte
		var templateData receiverTemplateDto

		templateData = receiverTemplateDto{
			Method: r.Method,
			Header: r.Header,
			Params: r.URL.Query(),
		}

		if r.Method == http.MethodPost {
			body, err = ioutil.ReadAll(r.Body)
			if err == nil && len(body) > 0 {
				switch r.Header.Get("Content-Type") {
				case "application/json":
					err = json.Unmarshal(body, &post)
					if err == nil {
						templateData.Body = post
					}
				case "application/xml":
					err = xml.Unmarshal(body, &post)
					if err == nil {
						templateData.Body = post
					}
				}
			}
		}

		buf := &bytes.Buffer{}
		tpl, err := template.New("test").Funcs(
			template.FuncMap{
				"getValue": func(key string, values map[string][]string) string {
					if values == nil {
						return ""
					}
					vs := values[key]
					if len(vs) == 0 {
						return ""
					}
					return vs[0]
				},
			},
		).Parse("**Пришли данные: -> {{ .Method }} имя {{ .Body.name }} и логин {{ .Body.login }}, а параметр {{  getValue \"test\" .Params }} **")

		if err == nil {
			err = tpl.Execute(buf, templateData)
		}

		channel, err := notifyChannel.New(team.Channel, a.processor.Tpl)
		_, err = channel.SendTeamMessage(notifyChannel.Message{Text: buf.String()})

		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		fmt.Println("Data receive!")
		fmt.Println("Data: ", templateData)
		fmt.Println("Result: ", buf.String())

		renderJson(w, http.StatusOK, ResponseSuccess{Result: "Ok"})
	})
}

func createReceiverResponseItem(receiver *domain.Receiver) receiverResponseItem {
	return receiverResponseItem{
		Id:       receiver.Id,
		Active:   receiver.Active,
		Type:     string(receiver.Type),
		Format:   string(receiver.Format),
		Settings: receiver.Settings,
	}
}

// swagger:route GET /team/{teamId}/receivers/{receiverId} Receivers ReceiversGet
//
// Информация о приемнике.
//
// Responses:
//        200: Receiver
//        404: ResponseError
func (a apiHandler) ReceiverGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		receiver := a.store.GetReceiver(vars["teamId"], vars["receiverId"])
		if receiver == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Task #%s on team #%s not found", vars["teamId"], vars["receiverId"])})
			return
		}

		renderJson(w, http.StatusOK, createReceiverResponseItem(receiver))
	})
}

// swagger:route GET /team/{teamId}/receivers Receivers ReceiversList
//
// Список приемников команды.
//
// Responses:
//        200: ReceiversList
func (a apiHandler) ReceiverList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		list := make([]receiverResponseItem, 0)
		for _, receiver := range a.store.GetReceivers(vars["teamId"]) {
			list = append(list, createReceiverResponseItem(receiver))
		}

		renderJson(w, http.StatusOK, receiverResponseList{Result: list})
	})
}

// swagger:route POST /team/{teamId}/receivers Receivers ReceiversAdd
//
// Создание приемника для команды.
//
// Responses:
//        200: Receiver
//        400: ResponseError
func (a apiHandler) ReceiverAdd() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var request *receiverResponseItem
		var receiver *domain.Receiver
		var body []byte
		var err error

		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			err = json.Unmarshal(body, &request)
		}

		receiver, err = domain.NewReceiver(
			request.Active,
			*domain.GetReceiverType(request.Type),
			request.Settings,
			*domain.GetReceiverFormat(request.Format),
		)
		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.AddReceiver(vars["teamId"], receiver)
		response := createReceiverResponseItem(receiver)
		renderJson(w, http.StatusOK, response)
	})
}

// swagger:route PATCH /team/{teamId}/receivers/{receiverId} Receivers ReceiversEdit
//
// Изменение приемника.
//
// Responses:
//        200: Receiver
//        404: ResponseError
//        400: ResponseError
func (a apiHandler) ReceiverEdit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var request *receiverResponseItem
		var receiver *domain.Receiver
		var body []byte
		var err error

		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			err = json.Unmarshal(body, &request)
		}

		receiver = a.store.GetReceiver(vars["teamId"], vars["receiverId"])
		if receiver == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Receiver #%s on team #%s not found", vars["teamId"], vars["receiverId"])})
			return
		}

		err = receiver.Edit(
			request.Active,
			*domain.GetReceiverType(request.Type),
			request.Settings,
			*domain.GetReceiverFormat(request.Format),
		)
		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.EditReceiver(vars["teamId"], receiver)
		response := createReceiverResponseItem(receiver)
		renderJson(w, http.StatusOK, response)
	})
}

// swagger:route DELETE /team/{teamId}/receivers/{receiverId} Receivers ReceiversDelete
//
// Удаление приемника.
//
// Responses:
//        200: ResponseSuccess
//        404: ResponseError
func (a apiHandler) ReceiverDelete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		receiver := a.store.GetReceiver(vars["teamId"], vars["receiverId"])
		if receiver == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Receiver #%s on team #%s not found", vars["teamId"], vars["receiverId"])})
			return
		}

		a.store.DeleteReceiver(vars["teamId"], receiver)

		renderJson(w, http.StatusOK, ResponseSuccess{Result: "ok"})
	})
}
