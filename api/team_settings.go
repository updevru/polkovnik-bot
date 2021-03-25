package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"polkovnik/domain"
)

type teamSettingsResponseItem struct {
	Id                     string            `json:"id"`
	Title                  string            `json:"title"`
	NotifyChannelType      string            `json:"notify_channel_type"`
	NotifyChannelChannelId string            `json:"notify_channel_channel_id"`
	NotifyChannelSettings  map[string]string `json:"notify_channel_settings"`
	IssueTrackerType       string            `json:"issue_tracker_type"`
	IssueTrackerSettings   map[string]string `json:"issue_tracker_settings"`
	MinWorkLog             int               `json:"min_work_log"`
}

func createTeamSettingsResponseItem(team *domain.Team) teamSettingsResponseItem {
	result := teamSettingsResponseItem{
		Id:         team.Id,
		Title:      team.Title,
		MinWorkLog: team.MinWorkLog,
	}

	if team.Channel != nil {
		result.NotifyChannelType = team.Channel.Type
		result.NotifyChannelChannelId = team.Channel.ChannelId
		result.NotifyChannelSettings = team.Channel.Settings
	}

	if team.IssueTracker != nil {
		result.IssueTrackerType = team.IssueTracker.Type
		result.IssueTrackerSettings = team.IssueTracker.Settings
	}

	return result
}

func (a apiHandler) TeamSettingsGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		team := a.store.GetTeam(vars["teamId"])
		if team == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Team #%s not found", vars["teamId"])})
			return
		}

		renderJson(w, http.StatusOK, createTeamSettingsResponseItem(team))
	})
}

func (a apiHandler) TeamSettingsEdit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var request *teamSettingsResponseItem
		var body []byte
		var err error

		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			err = json.Unmarshal(body, &request)
		}

		team := a.store.GetTeam(vars["teamId"])
		if team == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Team #%s not found", vars["teamId"])})
			return
		}

		err = team.EditSettings(
			request.Title,
			request.NotifyChannelType,
			request.NotifyChannelChannelId,
			request.NotifyChannelSettings,
			request.IssueTrackerType,
			request.IssueTrackerSettings,
			request.MinWorkLog,
		)

		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.EditTeam(team)

		renderJson(w, http.StatusOK, createTeamSettingsResponseItem(team))
	})
}
