package domain

type Config struct {
	Version float32
	Teams   []*Team
}

func NewConfig() *Config {
	return &Config{
		Version: 1,
		Teams:   []*Team{},
	}
}
