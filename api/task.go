package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type taskResponseItem struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	LastRunTime string `json:"last_run_time"`
}

type taskResponseList struct {
	Result []taskResponseItem `json:"result"`
}

func (a apiHandler) TaskList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var list []taskResponseItem
		for _, task := range a.store.GetTasks(vars["teamId"]) {
			row := taskResponseItem{
				Id:          task.Id,
				Type:        task.Type,
				LastRunTime: task.LastRunTime.Format("02-01-2006 15:04:05"),
			}
			list = append(list, row)
		}

		renderJson(w, http.StatusOK, taskResponseList{Result: list})
	})
}
