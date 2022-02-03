package repository

import (
	"polkovnik/adapter/storage"
	"polkovnik/domain"
	"testing"
)

func TestRepository_Flush(t *testing.T) {
	t.Run("Flush if has no updated", func(t *testing.T) {
		repo := &Repository{
			data: storage.NewConfigMemory(),
		}

		result, err := repo.Flush()

		if err != nil {
			t.Errorf("Flush() error = %v", err)
		}

		if result == true {
			t.Errorf("Flush() update if empty changes")
		}
	})

	t.Run("Flush if has updated", func(t *testing.T) {
		repo := &Repository{
			data: storage.NewConfigMemory(),
		}
		config := &domain.Config{}

		if repo.update(config) != true {
			t.Errorf("update() error, need true")
		}

		result, err := repo.Flush()

		if err != nil {
			t.Errorf("Flush() error = %v", err)
		}

		if result == false {
			t.Errorf("Flush() update if not empty changes")
		}
	})
}

func TestRepository_GetVersion(t *testing.T) {
	t.Run("Get version number", func(t *testing.T) {
		configWithVersion := &domain.Config{Version: float32(1)}
		repo := &Repository{config: configWithVersion}

		if repo.GetVersion() != configWithVersion.Version {
			t.Errorf("GetVersion() get version = %v, want version = %v", repo.GetVersion(), configWithVersion.Version)
		}
	})
}

func TestRepository_Load(t *testing.T) {
	t.Run("Loading config", func(t *testing.T) {
		repo := &Repository{
			data: storage.NewConfigMemory(),
		}

		err := repo.Load()

		if err != nil {
			t.Errorf("Load() error = %v", err)
		}
	})
}

func TestRepository_UpVersion(t *testing.T) {
	t.Run("Up version number", func(t *testing.T) {
		configWithVersion := &domain.Config{Version: float32(1)}

		var newVersion float32
		newVersion = 2

		repo := &Repository{
			data:   storage.NewConfigMemory(),
			config: configWithVersion,
		}
		repo.UpVersion(newVersion)
		if repo.GetVersion() != newVersion {
			t.Errorf("UpVersion() get version = %v, want version = %v", repo.GetVersion(), newVersion)
		}

		result, _ := repo.Flush()
		if result != true {
			t.Errorf("UpVersion() after up version repository must be flushed")
		}
	})
}

func TestRepository_update(t *testing.T) {
	t.Run("Update data after changes", func(t *testing.T) {

		repo := &Repository{data: storage.NewConfigMemory()}
		result := repo.update(&domain.Config{Version: 1})

		if result != true {
			t.Errorf("update() get result = %v, want result = %v", result, true)
		}

		result, _ = repo.Flush()
		if result != true {
			t.Errorf("update() after update repository must be flushed")
		}

		result, _ = repo.Flush()
		if result != false {
			t.Errorf("update() next flush not affected if empty changes")
		}
	})
}
