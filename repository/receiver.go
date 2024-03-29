package repository

import "polkovnik/domain"

func (r Repository) GetReceiver(teamId string, receiverId string) *domain.Receiver {
	team := r.GetTeam(teamId)
	for _, row := range team.Receivers {
		if row.Id == receiverId {
			return row
		}
	}

	return nil
}

func (r Repository) GetReceiverById(receiverId string) *domain.Receiver {
	for _, team := range r.GetTeams() {
		for _, row := range team.Receivers {
			if row.Id == receiverId {
				return row
			}
		}
	}

	return nil
}

func (r Repository) GetTeamByReceiver(receiverId string) *domain.Team {
	for _, team := range r.GetTeams() {
		for _, row := range team.Receivers {
			if row.Id == receiverId {
				return team
			}
		}
	}

	return nil
}

func (r Repository) GetReceivers(teamId string) []*domain.Receiver {
	var result []*domain.Receiver
	team := r.GetTeam(teamId)
	if team == nil {
		return result
	}

	for _, receiver := range team.Receivers {
		result = append(result, receiver)
	}

	return result
}

func (r *Repository) AddReceiver(teamId string, receiver *domain.Receiver) bool {
	team := r.GetTeam(teamId)
	team.Receivers = append(team.Receivers, receiver)

	return r.update(r.config)
}

func (r *Repository) EditReceiver(teamId string, receiver *domain.Receiver) bool {
	team := r.GetTeam(teamId)
	for i, row := range team.Receivers {
		if row.Id == receiver.Id {
			team.Receivers[i] = receiver
			return r.update(r.config)
		}
	}

	return false
}

func (r *Repository) DeleteReceiver(teamId string, receiver *domain.Receiver) bool {
	team := r.GetTeam(teamId)
	for index, row := range team.Receivers {
		if row.Id == receiver.Id {
			team.Receivers = append(team.Receivers[:index], team.Receivers[index+1:]...)
			return r.update(r.config)
		}
	}

	return false
}
