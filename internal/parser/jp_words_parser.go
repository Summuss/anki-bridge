package parser

import (
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/util"
)

type JPWordsParser struct {
}

func (w JPWordsParser) NoteName() string {
	return "Jp Words"
}

func (w JPWordsParser) Split(rowNotes string) ([]string, error) {
	return util.SplitByNoIndentLine(rowNotes)
}

func (w JPWordsParser) Check(note string) error {
	//TODO implement me
	panic("implement me")
}

func (w JPWordsParser) Parse(note string) (model.Model, error) {
	//TODO implement me
	panic("implement me")
}
