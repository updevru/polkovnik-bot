package storage

import "polkovnik/domain"

type ConfigInterface interface {
	Load() (*domain.Config, error)

	Update(*domain.Config) error
}
