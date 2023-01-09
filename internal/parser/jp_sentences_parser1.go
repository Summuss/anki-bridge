package parser

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
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

func (J JPSentencesParser1) Match(note string, noteType common.NoteInfo) bool {
	return common.NoteType_JPSentences_Name == noteType.Name && jpSentencesParser1Pattern.MatchString(note)
}

var jpSentencesParser1Pattern = regexp.MustCompile(`(?m)\A\s*^-\s*(?P<Sentence>\S+.*$)(?P<Words>(\n^\t-\s*.+$\n^\t\t-\s*.+$\n^\t\t-\s*.+$)+)\s*\z`)

func (J JPSentencesParser1) Check(note string, _ common.NoteInfo) error {
	notePreproc := common.PreprocessNote(note)
	if !jpSentencesParser1Pattern.MatchString(notePreproc) {
		return fmt.Errorf("note synatx error:\n%s", note)
	}
	submatches := jpSentencesParser1Pattern.FindStringSubmatch(notePreproc)
	wordsRaw := submatches[jpSentencesParser1Pattern.SubexpIndex("Words")]
	wordsRaw = common.UnIndent(wordsRaw)
	wordNoteInfo := *config.Conf.GetNoteInfoByName(common.NoteType_JPWords_Name)
	word_notes, err := splitter.Split(wordsRaw, wordNoteInfo)
	if err != nil {
		return fmt.Errorf("%s\n in note:\n%s", err.Error(), note)
	}
	errorList := lo.Map(
		*word_notes, func(item string, _ int) error {
			err = jpWordsParser.Check(item, wordNoteInfo)
			if err != nil {
				return fmt.Errorf("%s\nin note:\n%s", err.Error(), note)
			}
			return nil
		},
	)
	return common.MergeErrors(errorList)
}

func (J JPSentencesParser1) Parse(note string, noteType common.NoteInfo) (model.IModel, error) {
	notePreproc := common.PreprocessNote(note)
	r, _ := regexp.Compile(`(?m)\A\s*^-\s*(?P<Sentence>\S+.*$)(?P<Words>(\n^\t-\s*.+$\n^\t\t-\s*.+$\n^\t\t-\s*.+$)+)\s*\z`)
	submatches := r.FindStringSubmatch(notePreproc)
	sentence := submatches[r.SubexpIndex("Sentence")]
	wordsRaw := submatches[r.SubexpIndex("Words")]
	wordsRaw = common.UnIndent(wordsRaw)
	wordNoteInfo := *config.Conf.GetNoteInfoByName(common.NoteType_JPWords_Name)
	word_notes, _ := splitter.Split(wordsRaw, wordNoteInfo)
	var err error
	words := lo.Map(
		*word_notes, func(item string, _ int) *model.JPWord {
			word, e := jpWordsParser.Parse(item, wordNoteInfo)
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
