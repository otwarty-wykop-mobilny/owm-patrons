package common

import (
	"encoding/json"
	"io/ioutil"
)

// Config is a configuration object with non-sensitive data connected with the app's backend.
// Contains mainly GitHub related data.
type Config struct {
	BotUsername     string
	RepoAuthor      string
	RepoName        string
	Branch          string
	PatronsFilePath string
}

// Loads the configuration from json, panics on error
func LoadConfig() *Config {
	cfg := &Config{}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		panic(err)
	}
	return cfg
}
