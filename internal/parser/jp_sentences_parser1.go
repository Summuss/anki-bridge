package parser

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/splitter"
	"regexp"
)

var jpSentencesParser1 = JPSentencesParser1{}

func init() {
	*parsers = append(*parsers, jpSentencesParser1)
}

type JPSentencesParser1 struct {
	baseParser
}

func (J JPSentencesParser1) Match(note string, noteType common.NoteType) bool {
	return common.NoteType_JPSentences == noteType && jpSentencesParser1Pattern.MatchString(note)
}

var jpSentencesParser1Pattern = regexp.MustCompile(`(?m)\A\s*^-\s*(?P<Sentence>\S+.*$)(?P<Words>(\n^\t-\s*.+$\n^\t\t-\s*.+$\n^\t\t-\s*.+$)+)\s*\z`)

func (J JPSentencesParser1) Check(note string, _ common.NoteType) error {
	notePreproc := common.PreprocessNote(note)
	if !jpSentencesParser1Pattern.MatchString(notePreproc) {
		return fmt.Errorf("note synatx error:\n%s", note)
	}
	submatches := jpSentencesParser1Pattern.FindStringSubmatch(notePreproc)
	wordsRaw := submatches[jpSentencesParser1Pattern.SubexpIndex("Words")]
	wordsRaw = common.UnIndent(wordsRaw)
	word_notes, err := splitter.Split(wordsRaw, common.NoteType_JPWords)
	if err != nil {
		return fmt.Errorf("%s\n in note:\n%s", err.Error(), note)
	}
	errorList := lo.Map(
		*word_notes, func(item string, _ int) error {
			err = jpWordsParser.Check(item, "")
			if err != nil {
				return fmt.Errorf("%s\nin note:\n%s", err.Error(), note)
			}
			return nil
		},
	)
	return common.MergeErrors(errorList)
}

func (J JPSentencesParser1) Parse(note string, noteType common.NoteType) (model.IModel, error) {
	notePreproc := common.PreprocessNote(note)
	r, _ := regexp.Compile(`(?m)\A\s*^-\s*(?P<Sentence>\S+.*$)(?P<Words>(\n^\t-\s*.+$\n^\t\t-\s*.+$\n^\t\t-\s*.+$)+)\s*\z`)
	submatches := r.FindStringSubmatch(notePreproc)
	sentence := submatches[r.SubexpIndex("Sentence")]
	wordsRaw := submatches[r.SubexpIndex("Words")]
	wordsRaw = common.UnIndent(wordsRaw)
	word_notes, _ := splitter.Split(wordsRaw, common.NoteType_JPWords)
	var err error
	words := lo.Map(
		*word_notes, func(item string, _ int) *model.JPWord {
			word, e := jpWordsParser.Parse(item, "")
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
