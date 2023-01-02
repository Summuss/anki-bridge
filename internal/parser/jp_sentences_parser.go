package parser

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/util"
	"regexp"
)

var jpSentencesParser = JPSentencesParser{}

func init() {
	*parsers = append(*parsers, jpSentencesParser)
}

type JPSentencesParser struct {
}

func (J JPSentencesParser) Match(noteName string) bool {
	return "Jp Sentences" == noteName
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

func (J JPSentencesParser) Parse(note string) (model.IModel, error) {
	notePreproc := util.PreprocessNote(note)
	r, _ := regexp.Compile(`(?m)\A\s*^-\s*(?P<Sentence>\S+.*$)(?P<Words>(\n^\t-\s*.+$\n^\t\t-\s*.+$\n^\t\t-\s*.+$)+)\s*\z`)
	submatches := r.FindStringSubmatch(notePreproc)
	sentence := submatches[r.SubexpIndex("Sentence")]
	wordsRaw := submatches[r.SubexpIndex("Words")]
	wordsRaw = util.UnIndent(wordsRaw)
	word_notes, _ := jpWordsParser.Split(wordsRaw)
	var err error
	words := lo.Map(
		word_notes, func(item string, _ int) *model.JPWord {
			word, e := jpWordsParser.Parse(item)
			if e != nil {
				err = e
				return nil
			}
			return word.(*model.JPWord)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse sentence `%s`failed, error:\n%s", sentence, err.Error())
	}
	return &model.JPSentence{Sentence: sentence, JPWords: &words}, nil
}
