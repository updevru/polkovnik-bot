package api

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"polkovnik/domain"
	"polkovnik/job"
	"polkovnik/repository"
	"strings"
	"time"
)

type apiHandler struct {
	store     *repository.Repository
	history   *repository.HistoryRepository
	processor *job.Processor
}

func NewApiHandler(repository *repository.Repository, history *repository.HistoryRepository, processor *job.Processor) *apiHandler {
	return &apiHandler{
		store:     repository,
		history:   history,
		processor: processor,
	}
}

// Стандартное сообщение об ошибке.
//
// swagger:response ResponseError
type ResponseError struct {
	//Сообщение об ошибке
	Error string `json:"error"`
}

// Стандартное сообщение об успешном выполнении запроса.
//
// swagger:response ResponseSuccess
type ResponseSuccess struct {
	//Результат
	Result string `json:"result"`
}

type ResponseList struct {
	Result []interface{} `json:"result"`
}

// Выходные
// swagger:model Weekend
type weekendItem struct {
	// Выходные дни
	WeekDays []string `json:"week_days"`
	// Интервалы дат
	Intervals []weekendInterval `json:"intervals"`
}

func (u weekendItem) createIntervals() []domain.WeekendInterval {
	var intervals []domain.WeekendInterval
	for _, interval := range u.Intervals {
		intervals = append(intervals, interval.createInterval())
	}

	return intervals
}

// swagger:model WeekendInterval
type weekendInterval struct {
	// Дата начала выходного
	Start string `json:"start"`
	// Дата окончания выходного (включительно до конца дня).
	End string `json:"end"`
}

func (i weekendInterval) createInterval() domain.WeekendInterval {
	start, _ := time.ParseInLocation("02-01-2006", i.Start, time.Local)
	end, _ := time.ParseInLocation("02-01-2006", i.End, time.Local)

	return domain.WeekendInterval{
		Start: start,
		End:   end,
	}
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

type SpaHandler struct {
	StaticPath string
	IndexPath  string
	Files      embed.FS
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.StaticPath + r.URL.Path

	//Если запрашивается папка, то нужно смотреть на индексный файл
	if strings.HasSuffix(path, "/") {
		path = path + h.IndexPath
	}

	_, err := h.Files.Open(path)

	//Если запрашиваемого файла нет, отдаем главную страницу. Далее отработает роутинг в JavaScript
	if os.IsNotExist(err) {
		content, _ := h.Files.ReadFile(h.StaticPath + "/" + h.IndexPath)
		w.WriteHeader(http.StatusOK)
		w.Write(content)
		return
	}

	contentStatic, _ := fs.Sub(fs.FS(h.Files), h.StaticPath)
	http.FileServer(http.FS(contentStatic)).ServeHTTP(w, r)
}
