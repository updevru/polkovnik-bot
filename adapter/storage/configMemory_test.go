package storage

import (
	"polkovnik/domain"
	"reflect"
	"testing"
)

func TestConfigMemory_Load(t *testing.T) {
	t.Run("Load config", func(t *testing.T) {
		c := &ConfigMemory{}
		want := &domain.Config{}
		got, err := c.Load()
		if err != nil {
			t.Errorf("Load() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Load() got = %v, want %v", got, want)
		}
	})
}

func TestConfigMemory_Update(t *testing.T) {
	t.Run("Update config", func(t *testing.T) {
		config := &domain.Config{
			Version: 1,
		}
		c := &ConfigMemory{
			data: config,
		}

		newConfig := &domain.Config{
			Version: 2,
		}

		err := c.Update(newConfig)
		if err != nil {
			t.Errorf("Update() error = %v", err)
			return
		}
		if config.Version == newConfig.Version {
			t.Errorf("Update() not affected newVersion = %v, Version %v", newConfig.Version, config.Version)
		}
	})
}
