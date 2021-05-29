package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"polkovnik/domain"
)

type teamResponseItem struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type teamResponseList struct {
	Result []teamResponseItem `json:"result"`
}

func createTeamResponseItem(team *domain.Team) teamResponseItem {
	return teamResponseItem{
		Id:    team.Id,
		Title: team.Title,
	}
}

func (a apiHandler) TeamList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var list []teamResponseItem
		list = []teamResponseItem{}
		for _, team := range a.store.GetTeams() {
			list = append(list, createTeamResponseItem(team))
		}

		renderJson(w, http.StatusOK, teamResponseList{Result: list})
	})
}

func (a apiHandler) TeamAdd() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request *teamResponseItem
		var body []byte
		var err error
		var team *domain.Team

		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			err = json.Unmarshal(body, &request)
		}

		team, err = domain.NewTeam(request.Title)

		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.AddTeam(team)

		renderJson(w, http.StatusOK, createTeamResponseItem(team))
	})
}
