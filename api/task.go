package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"polkovnik/domain"
	"strings"
)

type taskResponseItem struct {
	Id               string   `json:"id"`
	Type             string   `json:"type"`
	ScheduleWeekdays []string `json:"schedule_weekdays"`
	ScheduleHour     int      `json:"schedule_hour"`
	ScheduleMinute   int      `json:"schedule_minute"`
	LastRunTime      string   `json:"last_run_time"`
	Projects         string   `json:"projects"`
	Message          string   `json:"message"`
	DateModify       string   `json:"check_date_modify"`
}

type taskResponseList struct {
	Result []taskResponseItem `json:"result"`
}

func createTaskResponseItem(task *domain.Task) taskResponseItem {
	return taskResponseItem{
		Id:               task.Id,
		Type:             task.Type,
		LastRunTime:      task.LastRunTime.Format("02-01-2006 15:04:05"),
		ScheduleWeekdays: task.Schedule.WeekDays,
		ScheduleHour:     task.Schedule.Hour,
		ScheduleMinute:   task.Schedule.Minute,
		Projects:         strings.Join(task.Projects, ","),
		Message:          task.Message,
		DateModify:       task.DateModify,
	}
}

func (a apiHandler) TaskGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		task := a.store.GetTask(vars["teamId"], vars["taskId"])
		if task == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Task #%s on team #%s not found", vars["teamId"], vars["taskId"])})
			return
		}

		renderJson(w, http.StatusOK, createTaskResponseItem(task))
	})
}

func (a apiHandler) TaskList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var list []taskResponseItem
		for _, task := range a.store.GetTasks(vars["teamId"]) {
			list = append(list, createTaskResponseItem(task))
		}

		renderJson(w, http.StatusOK, taskResponseList{Result: list})
	})
}

func (a apiHandler) TaskAdd() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var request *taskResponseItem
		var task *domain.Task
		var body []byte
		var err error

		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			err = json.Unmarshal(body, &request)
		}

		task, err = domain.NewTask(
			request.Type,
			request.ScheduleWeekdays,
			request.ScheduleHour,
			request.ScheduleMinute,
			strings.Split(request.Projects, ","),
			request.Message,
			request.DateModify,
		)
		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.AddTask(vars["teamId"], task)
		response := createTaskResponseItem(task)
		renderJson(w, http.StatusOK, response)
	})
}

func (a apiHandler) TaskEdit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var request *taskResponseItem
		var task *domain.Task
		var body []byte
		var err error

		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			err = json.Unmarshal(body, &request)
		}

		task = a.store.GetTask(vars["teamId"], vars["taskId"])
		if task == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Task #%s on team #%s not found", vars["teamId"], vars["taskId"])})
			return
		}

		err = task.Edit(
			request.Type,
			request.ScheduleWeekdays,
			request.ScheduleHour,
			request.ScheduleMinute,
			strings.Split(request.Projects, ","),
			request.Message,
			request.DateModify,
		)
		if err != nil {
			renderJson(w, http.StatusBadRequest, &ResponseError{Error: err.Error()})
			return
		}

		a.store.EditTask(vars["teamId"], task)
		response := createTaskResponseItem(task)
		renderJson(w, http.StatusOK, response)
	})
}

func (a apiHandler) TaskDelete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		task := a.store.GetTask(vars["teamId"], vars["taskId"])
		if task == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Task #%s on team #%s not found", vars["teamId"], vars["taskId"])})
			return
		}

		a.store.DeleteTask(vars["teamId"], task)

		renderJson(w, http.StatusOK, ResponseSuccess{Result: "ok"})
	})
}
