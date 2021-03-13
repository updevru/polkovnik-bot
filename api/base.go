package api

import (
	"encoding/json"
	"net/http"
	"polkovnik/repository"
)

type apiHandler struct {
	store *repository.Repository
}

func NewApiHandler(repository *repository.Repository) *apiHandler {
	return &apiHandler{
		store: repository,
	}
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseSuccess struct {
	Result string `json:"result"`
}

type ResponseList struct {
	Result []interface{} `json:"result"`
}

func renderJson(w http.ResponseWriter, status int, data interface{}) {
	response, err := json.Marshal(data)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Methods", "*")

	if err != nil {
		resultError, _ := json.Marshal(ResponseError{Error: err.Error()})
		w.WriteHeader(http.StatusBadGateway)
		w.Write(resultError)
	} else {
		w.WriteHeader(status)
		w.Write(response)
	}
}
