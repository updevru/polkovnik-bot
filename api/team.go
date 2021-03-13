package api

import (
	"net/http"
)

type teamResponseItem struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type teamResponseList struct {
	Result []teamResponseItem `json:"result"`
}

func (a apiHandler) TeamList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var list []teamResponseItem
		for _, team := range a.store.GetTeams() {
			row := teamResponseItem{
				Id:    team.Id,
				Title: team.Title,
			}
			list = append(list, row)
		}

		renderJson(w, http.StatusOK, teamResponseList{Result: list})
	})
}
