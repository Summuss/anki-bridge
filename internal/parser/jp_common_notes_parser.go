package parser

import (
	"fmt"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/model"
	"regexp"
	"strings"
)

var jpCommonNoteParser = JPCommonNoteParser{}

func init() {
	*parsers = append(*parsers, &jpCommonNoteParser)
}

type JPCommonNoteParser struct {
	baseParser
}

func (J JPCommonNoteParser) Match(note string, noteType *common.NoteInfo) bool {
	return noteType.Name == common.NoteType_JPCommonNotes_Name
}

func (J JPCommonNoteParser) MiddleParse(note string, noteType *common.NoteInfo) (
	model.IModel, error,
) {
	notePreproc := common.PreprocessNote(note)
	reg := regexp.MustCompile(`(?m)\A\s*^-\s*(?P<Front>.+)(?P<Back>($\n^\t- .+)+)\s*\z`)
	if !reg.MatchString(note) {
		return nil, fmt.Errorf("note synatx error:\n%s", note)
	}
	submatches := reg.FindStringSubmatch(notePreproc)
	front := submatches[reg.SubexpIndex("Front")]
	back := submatches[reg.SubexpIndex("Back")]
	back = common.UnIndent(back)
	back = strings.TrimSpace(back)

	return &model.JPCommonNote{Front: front, Back: back}, nil
}

func (J JPCommonNoteParser) PostParse(iModel model.IModel) (model.IModel, error) {
	return iModel, nil
}
