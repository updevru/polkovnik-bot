package api

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"polkovnik/domain"
	"polkovnik/repository"
	"time"
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

type weekendItem struct {
	WeekDays  []string          `json:"week_days"`
	Intervals []weekendInterval `json:"intervals"`
}

func (u weekendItem) createIntervals() []domain.WeekendInterval {
	var intervals []domain.WeekendInterval
	for _, interval := range u.Intervals {
		intervals = append(intervals, interval.createInterval())
	}

	return intervals
}

type weekendInterval struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (i weekendInterval) createInterval() domain.WeekendInterval {
	start, _ := time.Parse("02-01-2006", i.Start)
	end, _ := time.Parse("02-01-2006", i.End)

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
	_, err := h.Files.Open(h.StaticPath + r.URL.Path)

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
