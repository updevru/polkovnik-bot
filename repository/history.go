package repository

import (
	"errors"
	storm "github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"polkovnik/domain"
)

type HistoryRepository struct {
	db *storm.DB
}

func CreateHistoryRepository(path string) (*HistoryRepository, error) {
	db, err := storm.Open(path)

	if err == nil {
		err = db.Init(domain.History{})
	}

	if err != nil {
		return nil, errors.New("Repository error on open db: " + err.Error())
	}

	return &HistoryRepository{
		db: db,
	}, nil
}

func (h *HistoryRepository) Close() error {
	return h.db.Close()
}

func (h *HistoryRepository) New(history *domain.History) error {
	err := h.db.Save(history)
	if err != nil {
		return errors.New("Repository error on save: " + err.Error())
	}
	return nil
}

func (h *HistoryRepository) GetLastByTaskId(taskId string, limit int, offset int) ([]domain.History, error) {
	var result []domain.History

	query := h.db.Select(q.And(q.Eq("TaskId", taskId))).OrderBy("Date").Reverse().Limit(limit).Skip(offset)
	err := query.Find(&result)

	if err != nil && err != storm.ErrNotFound {
		return result, errors.New("Repository error on GetLastByTaskId: " + err.Error())
	}

	return result, nil
}

func (h HistoryRepository) GetCountByTaskId(taskId string) (int, error) {
	query := h.db.Select(q.And(q.Eq("TaskId", taskId)))
	result, err := query.Count(domain.History{})

	if err != nil {
		return 0, errors.New("Repository error on GetCountByTaskId: " + err.Error())
	}

	return result, nil
}
