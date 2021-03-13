package repository

import "polkovnik/domain"

func (r Repository) GetTasks(teamId string) []*domain.Task {
	var result []*domain.Task
	team := r.GetTeam(teamId)
	if team == nil {
		return result
	}

	for _, task := range team.Tasks {
		result = append(result, task)
	}

	return result
}
