package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var Conf *Config

func init() {
	yamlFile, err := os.ReadFile("conf.yml")
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
