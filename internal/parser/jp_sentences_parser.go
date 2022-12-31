package parser

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/util"
	"regexp"
)

var jpSentencesParser = JPSentencesParser{}

type JPSentencesParser struct {
}

func (J JPSentencesParser) NoteName() string {
	return "Jp Sentences"
}

func (J JPSentencesParser) Split(rawNotes string) ([]string, error) {
	return util.SplitByNoIndentLine(rawNotes)
}

func (J JPSentencesParser) Check(note string) error {
	notePreproc := util.PreprocessNote(note)
	r, _ := regexp.Compile(`(?m)\A\s*^-\s*(?P<Sentence>\S+.*$)(?P<Words>(\n^\t-\s*.+$\n^\t\t-\s*.+$\n^\t\t-\s*.+$)+)\s*\z`)
	if !r.MatchString(notePreproc) {
		return fmt.Errorf("note synatx error:\n%s", note)
	}
	submatches := r.FindStringSubmatch(notePreproc)
	wordsRaw := submatches[r.SubexpIndex("Words")]
	wordsRaw = util.UnIndent(wordsRaw)
	word_notes, err := jpWordsParser.Split(wordsRaw)
	if err != nil {
		return fmt.Errorf("%s\n in note:\n%s", err.Error(), note)
	}
	errorList := lo.Map(
		word_notes, func(item string, _ int) error {
			err = jpWordsParser.Check(item)
			if err != nil {
				return fmt.Errorf("%s\nin note:\n%s", err.Error(), note)
			}
			return nil
		},
	)
	return util.MergeErrors(errorList)
}

func (J JPSentencesParser) Parse(note string) (model.model, error) {
	//TODO implement me
	panic("implement me")
}
