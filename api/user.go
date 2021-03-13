package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"polkovnik/domain"
)

type userResponseItem struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	NickName string `json:"nickname"`
}

type userResponseList struct {
	Result []userResponseItem `json:"result"`
}

func createUserResponseItem(user *domain.User) userResponseItem {
	return userResponseItem{
		Id:       user.Id,
		Name:     user.Name,
		Login:    user.Login,
		NickName: user.NickName,
	}
}

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

		user, err = domain.NewUser(request.Name, request.Login, request.NickName)
		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.AddUser(vars["teamId"], user)
		response := createUserResponseItem(user)
		renderJson(w, http.StatusOK, response)
	})
}

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

		err = user.Edit(request.Name, request.Login, request.NickName)
		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.EditUser(vars["teamId"], user)
		response := createUserResponseItem(user)
		renderJson(w, http.StatusOK, response)
	})
}

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
