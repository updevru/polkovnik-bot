package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"polkovnik/domain"
)

type ConfigFile struct {
	file string
	data *domain.Config
}

func NewConfigFile(file string) *ConfigFile {
	return &ConfigFile{
		file: file,
	}
}

func (c *ConfigFile) Load() (*domain.Config, error) {
	file, err := os.OpenFile(c.file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(fileContents) > 0 {
		err = json.Unmarshal(fileContents, &c.data)
	} else {
		c.data = domain.NewConfig()
	}

	return c.data, err
}

func (c ConfigFile) Update(config *domain.Config) error {
	c.data = config
	fileContents, err := json.Marshal(c.data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.file, fileContents, 0644)
}
