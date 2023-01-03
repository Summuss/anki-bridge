package config

import (
	"github.com/summuss/anki-bridge/internal/common"
	"gopkg.in/yaml.v3"
	"log"
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
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Fatalf("parse yml failed: %v", err)
	}
	if len(Conf.TTScmd) == 0 {
		log.Fatalf("conf file error: tts-cmd is empty")
	}
	if !Conf.RealMode {
		log.Println("warning: test mode on")
	}

}

type Config struct {
	MongoConnectURL  string                     `yaml:"mongo-connect-url"`
	DBName           string                     `yaml:"db-name"`
	AnkiAPIURL       string                     `yaml:"anki-api-url"`
	DefaultInputFile string                     `yaml:"default-input-file"`
	TTScmd           []string                   `yaml:"tts-cmd"`
	RealMode         bool                       `yaml:"real-mode"`
	NoteType2Desk    map[common.NoteType]string `yaml:"note-type-2-desk"`
}
