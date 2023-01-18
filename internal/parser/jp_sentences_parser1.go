package parser

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
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

func (J JPSentencesParser1) Match(note string, noteType *common.NoteInfo) bool {
	return common.NoteType_JPSentences_Name == noteType.Name && jpSentencesParser1Pattern.MatchString(note)
}

var jpSentencesParser1Pattern = regexp.MustCompile(`(?m)\A\s*^-\s*(?P<Sentence>\S+.*$)(?P<Words>(\n^\t-\s*.+$\n^\t\t-\s*.+$\n^\t\t-\s*.+$)+)\s*\z`)

func (J JPSentencesParser1) MiddleParse(note string, noteType *common.NoteInfo) (
	model.IModel, error,
) {
	notePreproc := common.PreprocessNote(note)
	if !jpSentencesParser1Pattern.MatchString(notePreproc) {
		return nil, fmt.Errorf("note synatx error:\n%s", note)
	}
	submatches := jpSentencesParser1Pattern.FindStringSubmatch(notePreproc)
	sentence := submatches[jpSentencesParser1Pattern.SubexpIndex("Sentence")]
	wordsRaw := submatches[jpSentencesParser1Pattern.SubexpIndex("Words")]
	wordsRaw = common.UnIndent(wordsRaw)
	wordNoteInfo := config.Conf.GetNoteInfoByName(common.NoteType_JPWords_Name)
	word_notes, err := splitter.Split(wordsRaw, wordNoteInfo)
	if err != nil {
		return nil, fmt.Errorf("%s\n in note:\n%s", err.Error(), note)
	}
	var merr *multierror.Error
	words := lo.Map(
		*word_notes, func(item string, _ int) *model.JPWord {
			w, err := jpWordsParser.MiddleParse(item, wordNoteInfo)
			if err != nil {
				merr = multierror.Append(merr, fmt.Errorf("%s\nin note:\n%s", err.Error(), note))
				return nil
			}
			w.SetParser(jpWordsParser)
			w.SetNoteInfo(wordNoteInfo)
			return w.(*model.JPWord)
		},
	)
	if merr.ErrorOrNil() != nil {
		return nil, merr.ErrorOrNil()
	}

	return &model.JPSentence{Sentence: sentence, JPWords: &words}, nil

}

func (J JPSentencesParser1) PostParse(iModel model.IModel) (model.IModel, error) {
	jpsentence := iModel.(*model.JPSentence)
	var merr *multierror.Error
	ws := lo.Map(
		*jpsentence.JPWords, func(item *model.JPWord, _ int) *model.JPWord {
			w, err := jpWordsParser.PostParse(item)
			if err != nil {
				merr = multierror.Append(merr, err)
				return nil
			}
			return w.(*model.JPWord)
		},
	)
	if merr.ErrorOrNil() != nil {
		return nil, fmt.Errorf(
			"parse sentence `%s`failed, error:\n%s", jpsentence.Sentence, merr.ErrorOrNil(),
		)
	}
	jpsentence.JPWords = &ws
	return jpsentence, nil
}
