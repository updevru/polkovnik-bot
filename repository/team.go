package repository

import "polkovnik/domain"

func (r Repository) GetTeams() []*domain.Team {
	var result []*domain.Team
	for _, team := range r.config.Teams {
		result = append(result, team)
	}
	return result
}

func (r Repository) GetTeam(id string) *domain.Team {
	for _, team := range r.config.Teams {
		if team.Id == id {
			return team
		}
	}

	return nil
}
