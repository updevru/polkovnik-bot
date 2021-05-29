package repository

import "polkovnik/domain"

type Repository struct {
	config *domain.Config
}

func NewRepository(config *domain.Config) *Repository {
	return &Repository{
		config: config,
	}
}
