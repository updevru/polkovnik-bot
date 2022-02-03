package storage

import (
	"polkovnik/domain"
)

type ConfigMemory struct {
	data *domain.Config
}

func NewConfigMemory() *ConfigMemory {
	return &ConfigMemory{}
}

func (c *ConfigMemory) Load() (*domain.Config, error) {
	c.data = &domain.Config{}
	return c.data, nil
}

func (c *ConfigMemory) Update(config *domain.Config) error {
	c.data = config

	return nil
}
