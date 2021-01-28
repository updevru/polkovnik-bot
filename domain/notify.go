package domain

type NotifyChannel struct {
	Type      string
	ChannelId string
	Settings  map[string]string
}
