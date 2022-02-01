package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"polkovnik/domain"
	receiverProc "polkovnik/receiver"
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
	Id string `json:"id"`
	//Вкл/Выкл
	//Если включен, то принимает данные, иначе не принимает
	Active bool `json:"active"`
	//Название приемника
	Title string `json:"title"`
	//Тип приемника
	//send_team_message - отправка сообщения команде
	Type string `json:"type"`
	//Формат тела запроса (JSON, XML, AUTO)
	Format string `json:"format"`
	//Настройки приемника
	//Каждому типу приемника соответствуют свои настройки
	Settings map[string]string `json:"settings"`
	//Адрес приемника
	Url string `json:"url"`
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

func (a apiHandler) Receive() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		dump, _ := httputil.DumpRequest(r, true)
		fmt.Printf("Recive request: %s", string(dump))

		receiver := a.store.GetReceiverById(vars["receiverId"])
		if receiver == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Reciver #%s not found", vars["receiverId"])})
			return
		}

		team := a.store.GetTeamByReceiver(vars["receiverId"])
		if team == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Team not found")})
			return
		}

		if receiver.Active == false {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Reciver #%s is disabled", vars["receiverId"])})
			return
		}

		receiverProcessor := receiverProc.Processor{Tpl: a.processor.Tpl}
		err := receiverProcessor.Run(team, receiver, receiverProc.CreateTemplateDto(receiver, r))

		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		fmt.Println("Data receive!")

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
		Title:    receiver.Title,
		Url:      fmt.Sprintf("/receive/%s", receiver.Id),
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
			request.Title,
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
			request.Title,
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
