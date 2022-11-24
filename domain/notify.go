package domain

type NotifyChannel struct {
	Type      string
	ChannelId string
	Settings  map[string]string
}

func NewNotifyChannel() *NotifyChannel {
	return &NotifyChannel{
		Settings: make(map[string]string),
	}
}
