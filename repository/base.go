package repository

import (
	"polkovnik/adapter/storage"
	"polkovnik/domain"
)

type Repository struct {
	data    storage.ConfigInterface
	config  *domain.Config
	updated bool
}

func NewRepository(data storage.ConfigInterface) *Repository {
	return &Repository{
		data:    data,
		updated: false,
	}
}

func (r *Repository) Load() error {
	config, err := r.data.Load()

	if err == nil {
		r.config = config
	}

	return err
}

func (r *Repository) Flush() (bool, error) {
	if r.updated == false {
		return false, nil
	}

	err := r.data.Update(r.config)
	if err == nil {
		r.updated = false
	}

	return true, err
}

func (r *Repository) update(config *domain.Config) bool {
	r.updated = true
	return true
}

func (r Repository) GetVersion() float32 {
	return r.config.Version
}

func (r *Repository) UpVersion(version float32) bool {
	r.config.Version = version
	r.updated = true
	return true
}
