package repository

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"polkovnik/domain"
)

type HistoryRepository struct {
	db *bolt.DB
}

func CreateHistoryRepository(path string) (*HistoryRepository, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &HistoryRepository{
		db: db,
	}, nil
}

func (h *HistoryRepository) Close() error {
	return h.db.Close()
}

func (h *HistoryRepository) New(history *domain.History) error {
	return h.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(history.TaskId))
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(history)
		if err != nil {
			return err
		}
		return b.Put([]byte(history.Id), encoded)
	})
}

func (h *HistoryRepository) GetLastByTaskId(taskId string, limit int, offset int) ([]domain.History, error) {
	var result []domain.History
	err := h.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(taskId)).Cursor()

		num := limit
		skip := offset
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if skip > 0 {
				skip--
				continue
			}

			var item domain.History
			err := json.Unmarshal(v, &item)
			if err != nil {
				return err
			}

			result = append(result, item)

			num--
			if num == 0 {
				break
			}
		}

		return nil
	})

	return result, err
}

func (h HistoryRepository) GetCountByTaskId(taskId string) (int, error) {
	result := 0
	err := h.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskId))
		b.ForEach(func(k, v []byte) error {
			result++
			return nil
		})
		return nil
	})

	return result, err
}
