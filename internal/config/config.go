package config

import (
	"fmt"
	"github.com/summuss/anki-bridge/internal/common"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var Conf *Config

func init() {
	confPath := os.Getenv("ANKI_BRIDGE_CONF")
	if len(confPath) == 0 {
		confPath = "conf.yml"
	}
	yamlFile, err := os.ReadFile(confPath)
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
	if Conf.ResourceFolder != "" {
		stat, err := os.Stat(Conf.ResourceFolder)
		if os.IsNotExist(err) {
			log.Fatalf("resource-folder %s not found", Conf.ResourceFolder)
		}
		if !stat.IsDir() {
			log.Fatalf("resource-folder %s is't a dictionay", Conf.ResourceFolder)
		}
	}
	if !Conf.RealMode {
		log.Println("warning: test mode on")
	}

}

type Config struct {
	MongoConnectURL         string                                    `yaml:"mongo-connect-url"`
	DBName                  string                                    `yaml:"db-name"`
	AnkiAPIURL              string                                    `yaml:"anki-api-url"`
	DefaultInputFile        string                                    `yaml:"default-input-file"`
	TTScmd                  []string                                  `yaml:"tts-cmd"`
	BackupCmd               [][]string                                `yaml:"backup-cmd"`
	RealMode                bool                                      `yaml:"real-mode"`
	DisableDuplicationCheck bool                                      `yaml:"disable-duplication-check"`
	ResourceFolder          string                                    `yaml:"resource-folder"`
	NoteInfo                map[common.NoteTypeName]map[string]string `yaml:"note-info"`
	KindCopySuffix          []string                                  `yaml:"kindle-copy-suffix"`

	noteInfoCacheByName  map[common.NoteTypeName]*common.NoteInfo
	noteInfoCacheByTitle map[string]*common.NoteInfo
}

func (c *Config) GetNoteInfoByTitle(noteTitle string) (*common.NoteInfo, error) {
	if c.noteInfoCacheByTitle == nil {
		c.noteInfoCacheByTitle = make(map[string]*common.NoteInfo)
	}
	if _, ok := c.noteInfoCacheByTitle[noteTitle]; ok {
		return c.noteInfoCacheByTitle[noteTitle], nil
	}

	for k, v := range c.NoteInfo {
		if v["title"] == noteTitle {
			info := &common.NoteInfo{
				Name:          k,
				Title:         v["title"],
				Desk:          v["desk"],
				AnkiNoteModel: v["anki-note-model"],
			}
			c.noteInfoCacheByTitle[noteTitle] = info
			return info, nil
		}

	}
	return nil, fmt.Errorf("can't find NoteInfo from conf by noteTitle %s", noteTitle)
}
func (c *Config) GetNoteInfoByName(name common.NoteTypeName) *common.NoteInfo {
	if c.noteInfoCacheByName == nil {
		c.noteInfoCacheByName = make(map[common.NoteTypeName]*common.NoteInfo)
	}
	if _, ok := c.noteInfoCacheByName[name]; ok {
		return c.noteInfoCacheByName[name]
	}

	for k, v := range c.NoteInfo {
		if k == name {
			info := &common.NoteInfo{
				Name:          common.NoteTypeName(k),
				Title:         v["title"],
				Desk:          v["desk"],
				AnkiNoteModel: v["anki-note-model"],
			}
			c.noteInfoCacheByName[name] = info
			return info
		}

	}
	panic(fmt.Sprintf("can't find NoteInfo from conf by noteName %s", name))
}
