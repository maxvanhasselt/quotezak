package app

import (
	"fmt"
	"io/ioutil"

	"git.code-cloppers.com/max/quotezak/bot"
	"git.code-cloppers.com/max/quotezak/db"
	"gopkg.in/yaml.v2"
)

// Config is a struct that contains the configuration data for the application.
type Config struct {
	Port     int        `yaml:"port"`
	Database db.Config  `yaml:"database"`
	Bot      bot.Config `yaml:"irc"`
}

// FromFile reads a yaml config file and returns a Config object
func (cfg *Config) FromFile(configFile string) error {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, cfg)
}

// ToString returns a string representation of the Config object and its children.
func (cfg *Config) ToString() string {
	return fmt.Sprintf("port: %d\n database:\n%s", cfg.Port, cfg.Database.ToString())
}
