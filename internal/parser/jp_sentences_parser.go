package parser

import (
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/util"
)

type JPSentencesParser struct {
}

func (J JPSentencesParser) NoteName() string {
	return "Jp Sentences"
}

func (J JPSentencesParser) Split(rawNotes string) ([]string, error) {
	return util.SplitByNoIndentLine(rawNotes)
}

func (J JPSentencesParser) Check(note string) error {
	//TODO implement me
	panic("implement me")
}

func (J JPSentencesParser) Parse(note string) (model.Model, error) {
	//TODO implement me
	panic("implement me")
}
