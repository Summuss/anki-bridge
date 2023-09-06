package parser

import (
	"fmt"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/model"
	"regexp"
	"strings"
)

var jpListenMojiParser = JpListenMojiParser{}

func init() {
	*parsers = append(*parsers, &jpListenMojiParser)
}

type JpListenMojiParser struct {
	baseParser
}

func (J JpListenMojiParser) Match(note string, noteType *common.NoteInfo) bool {
	return noteType.Name == common.NoteType_JpListenMoji_Name
}

func (J JpListenMojiParser) MiddleParse(note string, noteType *common.NoteInfo) (
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

func (J JpListenMojiParser) PostParse(iModel model.IModel) (model.IModel, error) {
	return iModel, nil
}
