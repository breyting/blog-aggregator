package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	path := home + configFileName

	jsonFile, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(byteValue, &config)

	return config, nil
}

func (c Config) SetUser(name string) error {
	c.CurrentUserName = name
	json, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return nil
}
