package storage

import "teamBot/domain"

type ConfigInterface interface {
	Load() (*domain.Config, error)

	Update(*domain.Config) error
}
