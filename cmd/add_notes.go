package cmd

import (
	"fmt"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/parser"
	"github.com/summuss/anki-bridge/internal/render"
	"github.com/summuss/anki-bridge/internal/util"
)

func AddNotes(text string) error {
	err := parser.CheckInput(text)
	if err != nil {
		return fmt.Errorf("input check error:\n%s", err.Error())
	}
	ms, err := parser.Parse(text)
	if err != nil {
		return fmt.Errorf("parse error:\n %s", err.Error())
	}
	return util.DoParallel(
		ms, func(m *model.IModel) error {
			desc := (*m).Desc()
			err = (*m).Save(model.MongoClient, config.Conf.DBName)
			if err != nil {
				if _, ok := err.(model.ExistError); ok {
					fmt.Printf("warnning: %s already existed, skip", desc)
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
			_ = util.DoParallel(
				resources, func(r *model.Resource) error {
					err := anki.StoreMedia(r)
					if err != nil {
						fmt.Printf(
							"store %s to anki for %s failed, error:\n%s", r.Metadata.FileName, desc,
							err.Error(),
						)

					}
					return nil
				},
			)

			return nil
		},
	)

}
