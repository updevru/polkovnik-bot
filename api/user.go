package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"polkovnik/domain"
)

// swagger:parameters UserGet UserEdit UserDelete
type userId struct {
	// ID Пользователя
	// in: path
	// required: true
	UserId string `json:"userId"`
}

// Участник команды
// swagger:model User
type userResponseItem struct {
	Id string `json:"id"`
	//Имя пользователя
	Name string `json:"name"`
	//Логин пользователя в системе задач
	Login string `json:"login"`
	//Ник пользователя в система чата
	NickName string `json:"nickname"`
	//Рабочее время пользователя
	Weekend weekendItem `json:"weekend"`
	//Вкл/Выкл
	//Если включен, то участвует во всех задачах
	Active bool `json:"active"`
}

// swagger:response User
type userResponseItemWrapper struct {
	// in: body
	Body userResponseItem `json:"body"`
}

// swagger:parameters UserAdd UserEdit
type userRequestItemWrapper struct {
	// in: body
	Body userResponseItem `json:"body"`
}

// Список участников команды
// swagger:model Users
type userResponseList struct {
	Result []userResponseItem `json:"result"`
}

// swagger:response UserList
type userResponseListWrapper struct {
	// in: body
	Body userResponseList `json:"body"`
}

func createUserResponseItem(user *domain.User) userResponseItem {
	var intervals []weekendInterval
	for _, row := range user.Weekend.Intervals {
		intervals = append(intervals, weekendInterval{
			Start: row.Start.Format("02-01-2006"),
			End:   row.End.Format("02-01-2006"),
		})
	}
	return userResponseItem{
		Id:       user.Id,
		Name:     user.Name,
		Login:    user.Login,
		NickName: user.NickName,
		Weekend: weekendItem{
			WeekDays:  user.Weekend.WeekDays,
			Intervals: intervals,
		},
		Active: user.Active,
	}
}

// swagger:route GET /team/{teamId}/users/{userId} Users UserGet
//
// Информация об участнике команды.
//
// Responses:
//        200: User
//        404: ResponseError
func (a apiHandler) UserGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		user := a.store.GetUser(vars["teamId"], vars["userId"])
		if user == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("User #%s on team #%s not found", vars["teamId"], vars["userId"])})
			return
		}

		renderJson(w, http.StatusOK, createUserResponseItem(user))
	})
}

// swagger:route GET /team/{teamId}/users Users UserList
//
// Список участников команды.
//
// Responses:
//        200: UserList
func (a apiHandler) UserList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var list []userResponseItem
		for _, user := range a.store.GetUsers(vars["teamId"]) {
			row := createUserResponseItem(user)
			list = append(list, row)
		}

		renderJson(w, http.StatusOK, userResponseList{Result: list})
	})
}

// swagger:route POST /team/{teamId}/users Users UserAdd
//
// Создание нового участника
//
// Responses:
//        200: User
//        400: ResponseError
func (a apiHandler) UserAdd() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var request *userResponseItem
		var user *domain.User
		var body []byte
		var err error

		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			err = json.Unmarshal(body, &request)
		}

		user, err = domain.NewUser(request.Name, request.Login, request.NickName, request.Weekend.WeekDays, request.Weekend.createIntervals())
		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.AddUser(vars["teamId"], user)
		response := createUserResponseItem(user)
		renderJson(w, http.StatusOK, response)
	})
}

// swagger:route PATCH /team/{teamId}/users/{userId} Users UserEdit
//
// Изменение участника
//
// Responses:
//        200: User
//        400: ResponseError
//        404: ResponseError
func (a apiHandler) UserEdit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var request *userResponseItem
		var user *domain.User
		var body []byte
		var err error

		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			err = json.Unmarshal(body, &request)
		}

		user = a.store.GetUser(vars["teamId"], vars["userId"])
		if user == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("User #%s on team #%s not found", vars["teamId"], request.Id)})
			return
		}

		err = user.Edit(request.Name, request.Login, request.NickName, request.Active, request.Weekend.WeekDays, request.Weekend.createIntervals())
		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.EditUser(vars["teamId"], user)
		response := createUserResponseItem(user)
		renderJson(w, http.StatusOK, response)
	})
}

// swagger:route DELETE /team/{teamId}/users/{userId} Users UserDelete
//
// Удаление участника команды
//
// Responses:
//        200: ResponseSuccess
//        404: ResponseError
func (a apiHandler) UserDelete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		user := a.store.GetUser(vars["teamId"], vars["userId"])
		if user == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("User #%s on team #%s not found", vars["teamId"], vars["userId"])})
			return
		}

		a.store.DeleteUser(vars["teamId"], user)

		renderJson(w, http.StatusOK, ResponseSuccess{Result: "ok"})
	})
}
