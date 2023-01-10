package cmd

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/parser"
	"github.com/summuss/anki-bridge/internal/render"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var PreBackUpFlg bool

func init() {
	rootCmd.AddCommand(addNotesCMD)
	addNotesCMD.Flags().BoolVarP(
		&PreBackUpFlg,
		"pre-backup", "b",
		false, "backup db before add notes execute",
	)
}

var addNotesCMD = &cobra.Command{
	Use:  "add_notes",
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		decks, err := anki.GetAllDecks()
		if err != nil {
			return err
		}
		models, err := anki.GetAllAnkiModels()
		if err != nil {
			return err
		}
		for name := range config.Conf.NoteInfo {
			if !slices.Contains(common.NoteTypeNameList, name) {
				return fmt.Errorf("note type name of %s not exist", name)
			}
			noteInfo := config.Conf.GetNoteInfoByName(name)
			targetDesk := noteInfo.Desk
			if targetDesk == "" {
				return fmt.Errorf("[[%s]]'s target desk not specified", name)
			}
			if !slices.Contains(decks, targetDesk) {
				return fmt.Errorf("[[%s]]'s target desk %s not exist", name, targetDesk)
			}

			targetAnkiModel := noteInfo.AnkiNoteModel
			if targetAnkiModel == "" {
				return fmt.Errorf("[[%s]]'s target anki model not specified", name)
			}
			if !slices.Contains(models, targetAnkiModel) {
				return fmt.Errorf("[[%s]]'s target anki model %s not exist", name, targetAnkiModel)
			}
		}
		if PreBackUpFlg {
			err := backupDB()
			if err != nil {
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var inputPath string
		if len(args) == 0 {
			inputPath = config.Conf.DefaultInputFile
		} else {
			inputPath = args[0]
		}
		fs, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("open file %s failed, %s", inputPath, err.Error())
		}
		bs, err := io.ReadAll(fs)
		if err != nil {
			return fmt.Errorf("read file %s failed, %s", inputPath, err.Error())
		}
		return addNotes(string(bs))

	},
}

func addNotes(text string) error {
	text = common.RemoveExtraInfo(text)
	err := parser.CheckInput(text)
	if err != nil {
		return fmt.Errorf("input check error:\n%s", err.Error())
	}
	ms, err := parser.Parse(text)
	if err != nil {
		return fmt.Errorf("parse error:\n %s", err.Error())
	}
	var insertNumMu sync.Mutex
	var skipNumMu sync.Mutex
	var insertNum int
	var skipNum int
	err = common.DoParallel(
		ms, func(m *model.IModel) error {
			desc := (*m).Desc()
			err = (*m).Save(model.MongoClient, config.Conf.DBName)
			if err != nil {
				if _, ok := err.(model.ExistError); ok {
					log.Printf("warnning: %s already existed, skip", desc)
					skipNumMu.Lock()
					skipNum++
					skipNumMu.Unlock()
					return nil
				} else {
					return fmt.Errorf("save %s to db failed,error:\n%s", desc, err.Error())
				}
			}
			card, err := render.Render(*m)
			if err != nil {
				return fmt.Errorf("render %s failed,error:\n%s", desc, err.Error())
			}

			card.ModelID = (*m).GetID().Hex()
			card.Collection = (*m).CollectionName()
			err = anki.AddCard(card)
			if err != nil {
				return fmt.Errorf("add  %s to anki failed,error:\n%s", desc, err.Error())
			}
			(*m).SetAnkiNoteId(card.ID)
			err = (*m).Save(model.MongoClient, config.Conf.DBName)
			if err != nil {
				return fmt.Errorf("write note id back for %s failed,error:\n%s", desc, err.Error())
			}
			resources := (*m).GetResources()
			_ = common.DoParallel(
				resources, func(r *model.Resource) error {
					err := anki.StoreMedia(r)
					if err != nil {
						log.Printf(
							"store %s to anki for %s failed, error:\n%s", r.Metadata.FileName, desc,
							err.Error(),
						)

					}
					return nil
				},
			)
			insertNumMu.Lock()
			insertNum = insertNum + 1
			insertNumMu.Unlock()
			return nil
		},
	)
	log.Printf("insert/skip/total: %d/%d/%d\n", insertNum, skipNum, len(*ms))
	return err

}
func backupDB() error {
	if len(config.Conf.BackupCmd) == 0 {
		return fmt.Errorf("backup db failed: backup-cmd not found")

	}
	timeFormat := "2006-01-02#15:04:05"
	dictionaryName := time.Now().Format(timeFormat)
	for _, cmd := range config.Conf.BackupCmd {
		cmd = lo.Map(
			cmd, func(item string, _ int) string {
				return strings.ReplaceAll(item, "$1", dictionaryName)
			},
		)
		_, err := common.Exec(cmd[0], cmd[1:]...)
		if err != nil {
			return fmt.Errorf("backup db failed, %s", err.Error())
		}

	}
	log.Println("backup db successfully")
	return nil

	/*	backup_dest_in_container := "/data/db/bak/{dest}"

		cmd := "docker exec mongo mongodump  -d anki -u {MONGO_USER} -p {MONGO_PASSWD} --authenticationDatabase admin  -o {backup_dest_in_container}"
		run_remote_cmd(cmd)
		run_remote_cmd("mv /home/summus/docker-data/mongodb/bak/{dest} /home/summus/bak/mongo")
	*/
}
