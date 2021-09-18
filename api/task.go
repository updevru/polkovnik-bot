package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"polkovnik/domain"
	"time"
)

// swagger:parameters TaskGet TaskEdit TaskDelete TaskRun TaskHistory
type taskId struct {
	// ID Задачи
	// in: path
	// required: true
	TaskId string `json:"taskId"`
}

// Задача
// swagger:model Task
type taskResponseItem struct {
	Id               string            `json:"id"`
	Type             string            `json:"type"`
	ScheduleWeekdays []string          `json:"schedule_weekdays"`
	ScheduleHour     int               `json:"schedule_hour"`
	ScheduleMinute   int               `json:"schedule_minute"`
	LastRunTime      string            `json:"last_run_time"`
	Active           bool              `json:"active"`
	Settings         map[string]string `json:"settings"`
}

// swagger:response Task
type taskResponseItemWrapper struct {
	// in: body
	Body taskResponseItem `json:"body"`
}

// swagger:parameters TaskAdd TaskEdit
type taskRequestWrapper struct {
	// in: body
	Body taskResponseItem `json:"body"`
}

//Список задач
// swagger:model Tasks
type taskResponseList struct {
	Result []taskResponseItem `json:"result"`
}

// swagger:response TaskList
type taskResponseListWrapper struct {
	// in: body
	Body taskResponseList `json:"body"`
}

func createTaskResponseItem(task *domain.Task) taskResponseItem {
	return taskResponseItem{
		Id:               task.Id,
		Type:             task.Type,
		LastRunTime:      task.LastRunTime.Format("02-01-2006 15:04:05"),
		ScheduleWeekdays: task.Schedule.WeekDays,
		ScheduleHour:     task.Schedule.Hour,
		ScheduleMinute:   task.Schedule.Minute,
		Active:           task.Active,
		Settings:         task.TaskSettings,
	}
}

// swagger:route GET /team/{teamId}/tasks/{taskId} Tasks TaskGet
//
// Информация о задании.
//
// Responses:
//        200: Task
//        404: ResponseError
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

// swagger:route GET /team/{teamId}/tasks Tasks TaskList
//
// Список заданий команды.
//
// Responses:
//        200: TaskList
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

// swagger:route POST /team/{teamId}/tasks Tasks TaskAdd
//
// Создание задания для команды.
//
// Responses:
//        200: Task
//        400: ResponseError
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
			request.Settings,
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

// swagger:route PATCH /team/{teamId}/tasks/{taskId} Tasks TaskEdit
//
// Изменение задания.
//
// Responses:
//        200: Task
//        404: ResponseError
//        400: ResponseError
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
			request.Active,
			request.Settings,
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

// swagger:route DELETE /team/{teamId}/tasks/{taskId} Tasks TaskDelete
//
// Удаление задания.
//
// Responses:
//        200: ResponseSuccess
//        404: ResponseError
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

// swagger:route POST /team/{teamId}/tasks/{taskId}/run Tasks TaskRun
//
// Запуск задания вручную.
//
// Responses:
//        200: ResponseSuccess
//        404: ResponseError
func (a apiHandler) TaskRun() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		team := a.store.GetTeam(vars["teamId"])
		if team == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Team #%s not found", vars["teamId"])})
			return
		}

		task := a.store.GetTask(vars["teamId"], vars["taskId"])
		if task == nil {
			renderJson(w, http.StatusNotFound, &ResponseError{Error: fmt.Sprintf("Task #%s on team #%s not found", vars["teamId"], vars["taskId"])})
			return
		}

		a.processor.ScheduleTask(team, task, time.Now().In(time.Local))

		renderJson(w, http.StatusOK, ResponseSuccess{Result: "ok"})
	})
}
