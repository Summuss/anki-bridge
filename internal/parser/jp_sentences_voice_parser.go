package parser

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/splitter"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var jpSentencesVoiceParser = JPSentencesVoiceParser{}

func init() {
	*parsers = append(*parsers, jpSentencesVoiceParser)
}

var _ = `
- #FILENAME [2023-01-16][20-52-30].mp3 #FILENAME [2023-01-16][20-52-32].mp3
	- あれが[[#red]]==偏屈==じゃなくて何なんだ！
	- 偏屈#[[Jp Words]]
		- 怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）
		- へんくつ 1 1
(?m)\A\s*^- #FILENAME (?P<file>\S+$)(\n^\t- \s*(?P<sentence1>.*?)\s*$)?(?P<words>(\n^\t- \S+\s*#\[\[Jp Words]]$\n^\t\t- .*$\n^\t\t- .*$)*)\s*\z
`
var jpSentencesVoiceParserPattern = regexp.MustCompile(`(?m)\A\s*^- (?P<files>(#FILENAME \S+\s*)+)$(\n^\t- \s*(?P<sentence>.*?)\s*$)?(?P<words>(\n^\t- \S+\s*#\[\[Jp Words]]$\n^\t\t- .*$\n^\t\t- .*$)*)\s*\z`)

type JPSentencesVoiceParser struct {
	baseParser
}

type jPSentencesVoiceMiddleInfo struct {
	filePaths []string
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
	fileNamesPart := submatches[jpSentencesVoiceParserPattern.SubexpIndex("files")]
	fileNamesPart = strings.TrimSpace(fileNamesPart)
	fileNames := regexp.MustCompile(`\s*#FILENAME\s*`).Split(fileNamesPart, -1)
	filePaths, err := computeFileLocation(fileNames)
	if err != nil {
		return nil, err
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

	jpSentence.SetMiddleInfo(&jPSentencesVoiceMiddleInfo{filePaths: filePaths})
	return &jpSentence, nil
}

func computeFileLocation(fileNames []string) ([]string, error) {
	fileNames = lo.Compact(fileNames)
	var merr *multierror.Error
	res := lo.Map(
		fileNames, func(fileName string, _ int) string {
			var filePath string
			filePath1 := path.Join(config.Conf.ResourceFolder, fileName)
			filePath2 := path.Join(getProcessedVoiceLocation(), fileName)
			if common.FileExists(filePath1) {
				filePath = filePath1
			} else if common.FileExists(filePath2) {
				filePath = filePath2
			} else {
				merr = multierror.Append(merr, fmt.Errorf("resource file %s not exist", fileName))
			}
			return filePath
		},
	)
	return res, merr.ErrorOrNil()

}

func (J JPSentencesVoiceParser) PostParse(iModel model.IModel) (model.IModel, error) {
	jpSentence := iModel.(*model.JPSentence)
	middleInfo := jpSentence.GetMiddleInfo().(*jPSentencesVoiceMiddleInfo)
	resources, err := processFilePaths(middleInfo.filePaths)
	if err != nil {
		return nil, err
	}
	jpSentence.SetResources(resources)
	return jpSentence, nil
}

func moveVoiceFile(filePath string) error {
	if filepath.Dir(filePath) == config.Conf.ResourceFolder {
		location := getProcessedVoiceLocation()
		stat, err := os.Stat(location)
		if os.IsNotExist(err) {
			err := os.Mkdir(location, os.ModeDir)
			if err != nil {
				return fmt.Errorf(
					"create dictionary %s for processed voice failed, %s", location, err.Error(),
				)
			}
		} else {
			if !stat.IsDir() {
				return fmt.Errorf("%s is not a dir", location)
			}
		}

		filename := filepath.Base(filePath)
		err = common.MoveFile(filePath, path.Join(location, filename))
		if err != nil {
			return fmt.Errorf(
				"move file %s to %s failed, %s", filePath, config.Conf.ResourceFolder, err.Error(),
			)
		}
	}
	return nil
}

func processFilePaths(filePaths []string) (*[]model.Resource, error) {
	var merr *multierror.Error
	resources := lo.Map(
		filePaths, func(filePath string, _ int) model.Resource {
			data, err := os.ReadFile(filePath)
			if err != nil {
				merr = multierror.Append(
					merr, fmt.Errorf("read file %s failed, %s", filePath, err.Error()),
				)
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
			err = moveVoiceFile(filePath)
			if err != nil {
				log.Printf("warning: move %s failed, %s\n", filePath, err.Error())
			}
			return *resource
		},
	)
	return &resources, merr.ErrorOrNil()
}

func getProcessedVoiceLocation() string {
	return path.Join(config.Conf.ResourceFolder, "processed")
}
