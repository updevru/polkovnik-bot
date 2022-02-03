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

func (r *Repository) EditTeam(team *domain.Team) bool {
	for i, row := range r.config.Teams {
		if row.Id == team.Id {
			r.config.Teams[i] = team
			return true
		}
	}

	return r.update(r.config)
}

func (r *Repository) AddTeam(team *domain.Team) bool {
	r.config.Teams = append(r.config.Teams, team)

	return r.update(r.config)
}
