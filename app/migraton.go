package app

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"polkovnik/domain"
)

func Migrate(config *domain.Config) error {
	result := false
	if config.Version < 1 {
		upgradeV1(config)
		config.Version = 1
		result = true
	}

	if config.Version < 1.1 {
		upgradeV1Dot1(config)
		config.Version = 1.1
		result = true
	}

	if config.Version < 1.2 {
		upgradeV1Dot2(config)
		config.Version = 1.2
		result = true
	}

	if result == true {
		log.Info("Config migrated to version ", config.Version, " successful!")
	}

	return nil
}

func upgradeV1(config *domain.Config) {
	for _, team := range config.Teams {
		if len(team.Id) == 0 {
			team.Id = uuid.NewString()
		}

		for _, user := range team.Users {
			if len(user.Id) == 0 {
				user.Id = uuid.NewString()
			}
		}

		for _, task := range team.Tasks {
			if len(task.Id) == 0 {
				task.Id = uuid.NewString()
			}
		}
	}
}

func upgradeV1Dot1(config *domain.Config) {
	for _, team := range config.Teams {
		for _, user := range team.Users {
			if user.Active == false {
				user.Active = true
			}
		}
	}
}

func upgradeV1Dot2(config *domain.Config) {
	for _, team := range config.Teams {
		for _, task := range team.Tasks {
			if task.Active == false {
				task.Active = true
			}
		}
	}
}
