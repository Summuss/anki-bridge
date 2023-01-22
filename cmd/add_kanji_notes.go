package cmd

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/render"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func init() {
	rootCmd.AddCommand(addKanjiCMD)
}

var addKanjiCMD = &cobra.Command{
	Use:  "add_kanji",
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return addKanji()
	},
}

func addKanji() error {
	dao := model.GetDao(model.MongoClient, config.Conf.DBName, &model.Kanji{})
	res, err := dao.FindMany(bson.D{})
	if err != nil {
		return fmt.Errorf("query kanji failed, %s", err.Error())
	}
	var merr *multierror.Error
	for _, item := range *res {
		card, err := render.Render(item)
		if err != nil {
			merr = multierror.Append(
				merr, fmt.Errorf("render kanji '%s' failed, %s", item.Kanji, err.Error()),
			)
		}
		card.ModelID = item.ID.Hex()
		card.Collection = item.CollectionName()
		err = anki.AddCard(card)
		if err != nil {
			merr = multierror.Append(
				merr,
				fmt.Errorf("add kanji '%s' card to anki failed, %s", item.Kanji, err.Error()),
			)
		}
		item.SetAnkiNoteId(card.ID)
		err = item.Save(model.MongoClient, config.Conf.DBName)
		if err != nil {
			merr = multierror.Append(
				merr, fmt.Errorf("kanji '%s':rewrite db failed, %s", item.Kanji, err.Error()),
			)
		}
	}
	size := len(*res)
	if merr != nil {
		log.Printf("success/total: %d/%d\n", size-merr.Len(), size)
	} else {
		log.Printf("success/total: %d/%d\n", size, size)
	}
	return merr.ErrorOrNil()
}
