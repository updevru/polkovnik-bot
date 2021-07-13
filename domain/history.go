package domain

import (
	"github.com/google/uuid"
	"time"
)

type historyStatus int

const (
	historyStatusSuccess historyStatus = iota + 1
	historyStatusError
)

type History struct {
	Id     string
	Date   time.Time
	TaskId string
	Logs   []string
	Status historyStatus
}

func NewHistory(taskId string) *History {
	return &History{Id: uuid.NewString(), Date: time.Now().In(time.Local), TaskId: taskId}
}

func (h *History) AddLine(line string) {
	h.Logs = append(h.Logs, line)
}

func (h *History) SetSuccess() {
	h.Status = historyStatusSuccess
}

func (h *History) IsSuccess() bool {
	return h.Status == historyStatusSuccess
}

func (h *History) SetError() {
	h.Status = historyStatusError
}

func (h *History) IsError() bool {
	return h.Status == historyStatusError
}
