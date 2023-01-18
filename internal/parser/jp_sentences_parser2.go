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
	"strings"
)

var jpSentencesParser2 = JPSentencesParser2{}

func init() {
	*parsers = append(*parsers, jpSentencesParser2)
}

var _ = `
- 君の一生の思い出、 {{cloze しかと}} {{cloze  見届け}} たぞ。
	- 確と #[[Jp Words]]
		- 〔たしかに〕确实,明确,准确.
		- しかと　2　４
	- 見届ける#[[Jp Words]]
		- 看到，看准，看清；一直看到最后，结束，用眼看，确认。（最後まで見る。また、しっかり見る。）
		- みとどける　0　７
	- this is a cross line explanation
	  xxxx
		- 2 level
		  aaaa
	- this is explanation2
(?m)\A\s*^- (\S+.*$)(\n^\t- \S+\s*#\[\[Jp Words]]$\n^\t\t- .*$\n^\t\t- .*$)*(\n^\t- .*($\n^\t+ .*$)*(\n^\t{2,}.*$)*)*\s*\z

`
var jpSentencesParser2Pattern = regexp.MustCompile(`(?m)\A\s*^- #2 (?P<sentence>\S+.*$)(?P<words>(\n^\t- \S+\s*#\[\[Jp Words]]$\n^\t\t- .*$\n^\t\t- .*$)*)(?P<addition>(\n^\t- .*($\n^\t+ .*$)*(\n^\t{2,}.*$)*)*)\s*\z`)

type JPSentencesParser2 struct {
	baseParser
}

func (p JPSentencesParser2) Priority() int {
	return -1
}

func (J JPSentencesParser2) Match(note string, noteType *common.NoteInfo) bool {
	return common.NoteType_JPSentences_Name == noteType.Name && jpSentencesParser2Pattern.MatchString(note)
}

func (J JPSentencesParser2) MiddleParse(note string, noteType *common.NoteInfo) (
	model.IModel, error,
) {
	notePreproc := common.PreprocessNote(note)
	if !jpSentencesParser2Pattern.MatchString(notePreproc) {
		return nil, fmt.Errorf("note synatx error:\n%s", note)
	}
	submatches := jpSentencesParser2Pattern.FindStringSubmatch(notePreproc)
	sentence := submatches[jpSentencesParser2Pattern.SubexpIndex("sentence")]
	wordsRaw := submatches[jpSentencesParser2Pattern.SubexpIndex("words")]
	wordsRaw = common.UnIndent(wordsRaw)
	wordsRaw = strings.ReplaceAll(wordsRaw, "#[[Jp Words]]", "")
	additionRaw := submatches[jpSentencesParser2Pattern.SubexpIndex("addition")]
	additionRaw = common.UnIndent(additionRaw)
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

	return &model.JPSentence{
		Sentence: sentence, JPWords: &words, Addition: strings.TrimSpace(additionRaw),
	}, nil

}

func (J JPSentencesParser2) PostParse(iModel model.IModel) (model.IModel, error) {
	return jpSentencesParser1.PostParse(iModel)
}
