package parser

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/splitter"
	"os"
	"path"
	"regexp"
	"strings"
)

var jpSentencesVoiceParser = JPSentencesVoiceParser{}

func init() {
	*parsers = append(*parsers, jpSentencesVoiceParser)
}

var _ = `
- #FILENAME [2023-01-16][20-52-30].mp3
	- あれが[[#red]]==偏屈==じゃなくて何なんだ！
	- 偏屈#[[Jp Words]]
		- 怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）
		- へんくつ 1 1
(?m)\A\s*^- #FILENAME (?P<file>\S+$)(\n^\t- \s*(?P<sentence1>.*?)\s*$)?(?P<words>(\n^\t- \S+\s*#\[\[Jp Words]]$\n^\t\t- .*$\n^\t\t- .*$)*)\s*\z
`
var jpSentencesVoiceParserPattern = regexp.MustCompile(`(?m)\A\s*^- #FILENAME (?P<file>\S+$)(\n^\t- \s*(?P<sentence>.*?)\s*$)?(?P<words>(\n^\t- \S+\s*#\[\[Jp Words]]$\n^\t\t- .*$\n^\t\t- .*$)*)\s*\z`)

type JPSentencesVoiceParser struct {
	baseParser
}

type jPSentencesVoiceMiddleInfo struct {
	filePath string
}

func (J JPSentencesVoiceParser) Match(note string, noteType *common.NoteInfo) bool {
	return common.NoteType_JPSentences_Voice_Name == noteType.Name
}

func (J JPSentencesVoiceParser) MiddleParse(note string, noteType *common.NoteInfo) (
	model.IModel, error,
) {
	notePreproc := common.PreprocessNote(note)
	if !jpSentencesVoiceParserPattern.MatchString(notePreproc) {
		return nil, fmt.Errorf("note synatx error:\n%s", note)
	}
	submatches := jpSentencesVoiceParserPattern.FindStringSubmatch(notePreproc)

	sentence := submatches[jpSentencesVoiceParserPattern.SubexpIndex("sentence")]
	fileName := submatches[jpSentencesVoiceParserPattern.SubexpIndex("file")]
	filePath := path.Join(config.Conf.ResourceFolder, fileName)
	if !common.FileExists(filePath) {
		return nil, fmt.Errorf("resource file %s not exist", filePath)
	}

	wordsRaw := submatches[jpSentencesVoiceParserPattern.SubexpIndex("words")]
	wordsRaw = common.UnIndent(wordsRaw)
	wordsRaw = strings.ReplaceAll(wordsRaw, "#[[Jp Words]]", "")

	wordNoteInfo := config.Conf.GetNoteInfoByName(common.NoteType_JPWords_Name)

	//check words
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

	jpSentence := model.JPSentence{
		Sentence: sentence, JPWords: &words,
	}

	jpSentence.SetMiddleInfo(&jPSentencesVoiceMiddleInfo{filePath: filePath})
	return &jpSentence, nil
}

func (J JPSentencesVoiceParser) PostParse(iModel model.IModel) (model.IModel, error) {
	jpSentence := iModel.(*model.JPSentence)
	middleInfo := jpSentence.GetMiddleInfo().(*jPSentencesVoiceMiddleInfo)
	filePath := middleInfo.filePath
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file %s failed, %s", filePath, err.Error())
	}
	filename := path.Base(filePath)
	filename = strings.Replace(filename, "[", "(", -1)
	filename = strings.Replace(filename, "]", ")", -1)
	resource := &model.Resource{
		Metadata: model.ResourceMetadata{
			FileName: filename, ResourceType: model.Sound, ExtName: ".mp3",
		},
	}
	resource.SetData(data)

	jpSentence.SetResources(&[]model.Resource{*resource})

	return jpSentence, nil
}
