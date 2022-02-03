package app

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"polkovnik/domain"
	"polkovnik/repository"
	"strings"
)

func Migrate(repository *repository.Repository) error {
	result := false
	if repository.GetVersion() < 1 {
		upgradeV1(repository)
		repository.UpVersion(1)
		result = true
	}

	if repository.GetVersion() < 1.1 {
		upgradeV1Dot1(repository)
		repository.UpVersion(1)
		result = true
	}

	if repository.GetVersion() < 1.2 {
		upgradeV1Dot2(repository)
		repository.UpVersion(1.2)
		result = true
	}

	if repository.GetVersion() < 1.3 {
		upgradeV1Dot3(repository)
		repository.UpVersion(1.3)
		result = true
	}

	if result == true {
		log.Info("Config migrated to version ", repository.GetVersion(), " successful!")
	}

	return nil
}

func upgradeV1(repository *repository.Repository) {
	for _, team := range repository.GetTeams() {
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

func upgradeV1Dot1(repository *repository.Repository) {
	for _, team := range repository.GetTeams() {
		for _, user := range team.Users {
			if user.Active == false {
				user.Active = true
			}
		}
	}
}

func upgradeV1Dot2(repository *repository.Repository) {
	for _, team := range repository.GetTeams() {
		for _, task := range team.Tasks {
			if task.Active == false {
				task.Active = true
			}
		}
	}
}

func upgradeV1Dot3(repository *repository.Repository) {
	for _, team := range repository.GetTeams() {
		for _, task := range team.Tasks {
			if len(task.TaskSettings) > 0 {
				continue
			}

			log.Info("Need update task ", task.Id, " on team ", team.Title)
			switch task.Type {
			case domain.CheckTeamWorkLog:
				settings := domain.TaskSettings{}
				settings["projects"] = strings.Join(task.Projects, ",")

				if len(task.DateModify) > 0 {
					settings["date_modify"] = task.DateModify
				}

				task.TaskSettings = settings
				break
			case domain.SendTeamMessage:
				settings := domain.TaskSettings{}
				settings["message"] = task.Message
				task.TaskSettings = settings
				break
			case domain.CheckUserWeekend:
				settings := domain.TaskSettings{}
				settings["date_modify"] = task.DateModify
				task.TaskSettings = settings
			}
		}
	}
}
