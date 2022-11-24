package domain

type IssueTracker struct {
	Type     string
	Settings map[string]string
}

func NewIssueTracker() *IssueTracker {
	return &IssueTracker{
		Settings: make(map[string]string),
	}
}
