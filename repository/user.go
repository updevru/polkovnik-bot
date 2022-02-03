package repository

import "polkovnik/domain"

func (r Repository) GetUser(teamId string, userId string) *domain.User {
	team := r.GetTeam(teamId)
	for _, row := range team.Users {
		if row.Id == userId {
			return row
		}
	}

	return nil
}

func (r Repository) GetUsers(teamId string) []*domain.User {
	var result []*domain.User
	team := r.GetTeam(teamId)
	if team == nil {
		return result
	}

	for _, user := range team.Users {
		result = append(result, user)
	}

	return result
}

func (r *Repository) AddUser(teamId string, user *domain.User) bool {
	team := r.GetTeam(teamId)
	team.Users = append(team.Users, user)

	return r.update(r.config)
}

func (r *Repository) EditUser(teamId string, user *domain.User) bool {
	team := r.GetTeam(teamId)
	for i, row := range team.Users {
		if row.Id == user.Id {
			team.Users[i] = user
			return r.update(r.config)
		}
	}

	return false
}

func (r *Repository) DeleteUser(teamId string, user *domain.User) bool {
	team := r.GetTeam(teamId)
	for index, row := range team.Users {
		if row.Id == user.Id {
			team.Users = append(team.Users[:index], team.Users[index+1:]...)
			return r.update(r.config)
		}
	}

	return false
}
