package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"polkovnik/domain"
)

// swagger:parameters MessageSend TaskGet TaskAdd TaskList TaskEdit TaskDelete TaskRun TaskHistory TeamsSettingsGet UserGet UserList UserAdd UserEdit UserDelete
type teamId struct {
	// ID команды
	// in: path
	// required: true
	TeamId string `json:"teamId"`
}

// swagger:model Team
type teamResponseItem struct {
	// ID команды
	Id string `json:"id"`
	// Название команды
	Title string `json:"title"`
}

// swagger:model Teams
type teamResponseList struct {
	Result []teamResponseItem `json:"result"`
}

// swagger:response TeamItem
type teamResponseWrapper struct {
	//in: body
	Body teamResponseItem `json:"body"'`
}

// swagger:parameters TeamAdd
type teamRequestWrapper struct {
	//in: body
	Body teamResponseItem `json:"body"'`
}

// swagger:response TeamList
type teamResponseListWrapper struct {
	//in: body
	Body teamResponseList `json:"body"'`
}

func createTeamResponseItem(team *domain.Team) teamResponseItem {
	return teamResponseItem{
		Id:    team.Id,
		Title: team.Title,
	}
}

// swagger:route GET /team Teams TeamList
//
// Список команд.
//
// Responses:
//        200: TeamList
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

// swagger:route POST /team Teams TeamAdd
//
// Создание команды.
//
// Responses:
//        200: TeamItem
//        400: ResponseError
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
