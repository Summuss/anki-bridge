package parser

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/splitter"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
)

type iParser interface {
	Match(note string, noteType *common.NoteInfo) bool
	MiddleParse(note string, noteType *common.NoteInfo) (model.IModel, error)
	PostParse(iModel model.IModel) (model.IModel, error)
	Priority() int
}

type baseParser struct{}

func (p baseParser) Priority() int {
	return 0
}

var (
	parsers = &[]iParser{}
)

func splitByNoteType(content string) (map[*common.NoteInfo]string, error) {
	content = content + "\n"
	r2, _ := regexp.Compile(`(?m)^\S+.*$\n`)
	splits := r2.Split(content, -1)
	if len(strings.TrimSpace(splits[0])) > 0 {
		return nil, fmt.Errorf("syntax error near %s", splits[0])
	}
	matches := r2.FindAllString(content, -1)
	noteTypeTitles := lo.Map(
		matches, func(item string, _ int) string {
			noteTypeTitle := strings.TrimSpace(item)
			noteTypeTitle = computeRealNoteName(noteTypeTitle)
			return noteTypeTitle
		},
	)

	noteTypeName2SubTxt := make(map[*common.NoteInfo]string)
	for i, noteTypeTitle := range noteTypeTitles {
		note := strings.TrimSpace(common.UnIndent(splits[i+1]))

		if len(noteTypeTitle) == 0 || len(note) == 0 {
			continue
		}
		noteInfo, err := config.Conf.GetNoteInfoByTitle(noteTypeTitle)
		if err != nil {
			return nil, err
		}
		noteTypeName2SubTxt[noteInfo] += "\n" + note
	}
	return noteTypeName2SubTxt, nil
}

func MiddleParse(txt string) (*[]model.IModel, error) {
	r, _ := regexp.Compile(`(?m)\A(\s*^\S+.*(\n^\t+.*$)*)*\s*\z`)
	if !r.MatchString(txt) {
		return nil, errors.New("input structure error")
	}
	noteTypeInfo2SubTxt, err := splitByNoteType(txt)
	if err != nil {
		return nil, err
	}
	size := len(noteTypeInfo2SubTxt)
	if size == 0 {
		return nil, nil
	}
	errCh := make(chan error, size)
	modelCh := make(chan model.IModel, 20)
	for noteInfo, subTxt := range noteTypeInfo2SubTxt {
		noteInfo := noteInfo
		subTxt := subTxt
		go func() {
			notes, err := splitter.Split(subTxt, noteInfo)
			if err != nil {
				errCh <- err
				return
			}

			err = common.DoParallel(
				notes, func(note *string) error {
					parser, err := findParser(noteInfo, *note)
					if err != nil {
						return err
					}
					m, err := parser.MiddleParse(*note, noteInfo)
					if err != nil {
						return err
					}

					m.SetNoteInfo(noteInfo)
					m.SetParser(parser)
					m.SetNoteTypeName(noteInfo.Name)
					modelCh <- m
					return nil
				},
			)
			errCh <- err
		}()
	}
	var errList []error
	var modelList []model.IModel
	for len(errList) != size {
		select {
		case err := <-errCh:
			errList = append(errList, err)
		case m := <-modelCh:
			modelList = append(modelList, m)
		}
	}
	close(modelCh)
	modelList = append(modelList, lo.ChannelToSlice(modelCh)...)
	return &modelList, common.MergeErrors(errList)
}

func FinalParse(middleModelList *[]model.IModel) (*[]model.IModel, error) {
	size := len(*middleModelList)
	modelCh := make(chan model.IModel, size)
	var merr *multierror.Error

	_ = common.DoParallel(
		middleModelList, func(m *model.IModel) error {
			parser := (*m).GetParser().(iParser)
			m2, err := parser.PostParse(*m)
			merr = multierror.Append(merr, err)
			modelCh <- m2
			return nil
		},
	)
	var ms []model.IModel
	for i := 0; i < size; i++ {
		ms = append(ms, <-modelCh)
	}
	if merr.ErrorOrNil() != nil {
		return nil, merr.ErrorOrNil()
	}
	return &ms, nil

}

func findParser(noteInfo *common.NoteInfo, note string) (iParser, error) {
	parserFlt := lo.Filter(
		*parsers, func(item iParser, index int) bool {
			return item.Match(note, noteInfo)
		},
	)
	if len(parserFlt) < 1 {
		return nil, fmt.Errorf(
			"can't found parser for note type of %v and note:\n%s", noteInfo, note,
		)
	} else if len(parserFlt) > 1 {
		slices.SortFunc(
			parserFlt, func(a, b iParser) bool {
				return a.Priority() > b.Priority()
			},
		)
		if parserFlt[0].Priority() == parserFlt[1].Priority() {
			return nil, fmt.Errorf(
				"found multiple parser with same priority for note type of %v", noteInfo,
			)
		}
		return parserFlt[0], nil
	}
	p := parserFlt[0]
	return p, nil
}
func computeRealNoteName(rowNoteName string) string {
	r, _ := regexp.Compile(`^- \[\[(.*)]]\s*$`)
	if !r.MatchString(rowNoteName) {
		return rowNoteName
	}
	return r.FindStringSubmatch(rowNoteName)[1]
}
