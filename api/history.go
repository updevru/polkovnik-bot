package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"polkovnik/domain"
	"strconv"
)

type taskHistoryResponseItem struct {
	Id      string   `json:"id"`
	Date    string   `json:"date"`
	Logs    []string `json:"logs"`
	Success bool     `json:"success"`
	Error   bool     `json:"error"`
}

type taskHistoryResponseList struct {
	Result []taskHistoryResponseItem `json:"result"`
	Total  int                       `json:"total"`
}

func createTaskHistoryResponseItem(row domain.History) taskHistoryResponseItem {
	return taskHistoryResponseItem{
		Id:      row.Id,
		Date:    row.Date.Format("02-01-2006 15:04:05"),
		Logs:    row.Logs,
		Success: row.IsSuccess(),
		Error:   row.IsError(),
	}
}

func (a apiHandler) HistoryList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var list []taskHistoryResponseItem
		list = []taskHistoryResponseItem{}

		size, _ := strconv.Atoi(r.URL.Query().Get("size"))
		if size == 0 {
			size = 30
		}
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}

		rows, err := a.history.GetLastByTaskId(vars["taskId"], size, (page-1)*size)

		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		for _, row := range rows {
			list = append(list, createTaskHistoryResponseItem(row))
		}

		total, _ := a.history.GetCountByTaskId(vars["taskId"])

		renderJson(w, http.StatusOK, taskHistoryResponseList{Result: list, Total: total})
	})
}
