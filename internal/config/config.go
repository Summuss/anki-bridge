package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var Conf *Config

func init() {
	path := os.Getenv("ANKI_BRIDGE_CONF")
	if len(path) == 0 {
		path = "conf.yml"
	}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("yamlFile.Get err   #%v ", err))
	}

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		panic(fmt.Errorf("parse yml failed: %v", err))
	}

}

type Config struct {
	MongoConnectURL string `yaml:"mongo-connect-url"`
	DBName          string `yaml:"db-name"`
	AnkiAPIURL      string `yaml:"anki-api-url"`
}
