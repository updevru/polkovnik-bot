package domain

type Team struct {
	Title        string
	Users        []*User
	Tasks        []*Task
	Channel      NotifyChannel
	Weekend      Weekend
	IssueTracker IssueTracker
	MinWorkLog   int
}

func (t *Team) AddUserPoint(user User, point int) bool {
	for _, item := range t.Users {
		if item.Login == user.Login {
			item.AddPoint(point)
			return true
		}
	}

	return false
}

func (t *Team) DeleteUserPoint(user User, point int) bool {
	for _, item := range t.Users {
		if item.Login == user.Login {
			item.DeletePoint(point)
			return true
		}
	}

	return false
}
