package repository

import "polkovnik/domain"

func (r Repository) GetTask(teamId string, taskId string) *domain.Task {
	team := r.GetTeam(teamId)
	for _, row := range team.Tasks {
		if row.Id == taskId {
			return row
		}
	}

	return nil
}

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

func (r Repository) AddTask(teamId string, task *domain.Task) bool {
	team := r.GetTeam(teamId)
	team.Tasks = append(team.Tasks, task)

	return true
}

func (r Repository) EditTask(teamId string, task *domain.Task) bool {
	team := r.GetTeam(teamId)
	for i, row := range team.Tasks {
		if row.Id == task.Id {
			team.Tasks[i] = task
			return true
		}
	}

	return false
}

func (r Repository) DeleteTask(teamId string, task *domain.Task) bool {
	team := r.GetTeam(teamId)
	for index, row := range team.Tasks {
		if row.Id == task.Id {
			team.Tasks = append(team.Tasks[:index], team.Tasks[index+1:]...)
			return true
		}
	}

	return false
}
