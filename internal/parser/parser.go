package parser

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/splitter"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
	"time"
)

type iParser interface {
	Match(note string, noteType common.NoteInfo) bool
	Check(note string, noteType common.NoteInfo) error
	Parse(note string, noteType common.NoteInfo) (model.IModel, error)
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

func CheckInput(txt string) error {
	r, _ := regexp.Compile(`(?m)\A(\s*^\S+.*(\n^\t+.*$)*)*\s*\z`)
	if !r.MatchString(txt) {
		return errors.New("input structure error")
	}
	noteTypeInfo2SubTxt, err := splitByNoteType(txt)
	if err != nil {
		return err
	}
	size := len(noteTypeInfo2SubTxt)
	if size == 0 {
		return nil
	}
	errCh := make(chan error, size)
	for noteInfo, subTxt := range noteTypeInfo2SubTxt {
		noteInfo := noteInfo
		subTxt := subTxt
		go func() {
			notes, err := splitter.Split(subTxt, *noteInfo)
			if err != nil {
				errCh <- err
				return
			}

			err = common.DoParallel(
				notes, func(note *string) error {
					parser, err := findParser(*noteInfo, *note)
					if err != nil {
						return err
					}
					err = parser.Check(*note, *noteInfo)
					if err != nil {
						return err
					}
					return nil
				},
			)
			errCh <- err
		}()
	}
	var errList []error
	for i := 0; i < size; i++ {
		errList = append(errList, <-errCh)
	}
	return common.MergeErrors(errList)
}

func Parse(text string) (*[]model.IModel, error) {
	var res []model.IModel
	noteTypeInfo2SubTxt, _ := splitByNoteType(text)
	size := len(noteTypeInfo2SubTxt)
	if size == 0 {
		return &res, nil
	}

	modelCh := make(chan model.IModel, 10)
	errCh := make(chan error, size)
	countCh := make(chan int, size)

	for noteInfo, subTxt := range noteTypeInfo2SubTxt {
		noteInfo := noteInfo
		subTxt := subTxt
		go func() {
			defer func() {
				countCh <- 1
			}()
			notes, _ := splitter.Split(subTxt, *noteInfo)
			err := common.DoParallel(
				notes, func(note *string) error {
					parser, _ := findParser(*noteInfo, *note)
					m, err := parser.Parse(*note, *noteInfo)
					if err != nil {
						return err
					}
					m.SetNoteTypeName(noteInfo.Name)
					modelCh <- m
					return nil
				},
			)
			if err != nil {
				errCh <- err
				return
			}
		}()
	}
	var count int
	for {
		if count == size {
			close(errCh)
			close(modelCh)
			break
		}
		select {
		case m := <-modelCh:
			res = append(res, m)
		case <-countCh:
			count++
		default:
			time.Sleep(20 * time.Millisecond)
		}
	}
	for m := range modelCh {
		res = append(res, m)
	}

	err := common.MergeErrors(lo.ChannelToSlice(errCh))
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func findParser(noteInfo common.NoteInfo, note string) (iParser, error) {
	parserFlt := lo.Filter(
		*parsers, func(item iParser, index int) bool {
			return item.Match(note, noteInfo)
		},
	)
	if len(parserFlt) < 1 {
		return nil, fmt.Errorf(
			"can't found parser for note type of %s and note:\n%s", noteInfo, note,
		)
	} else if len(parserFlt) > 1 {
		slices.SortFunc(
			parserFlt, func(a, b iParser) bool {
				return a.Priority() > b.Priority()
			},
		)
		if parserFlt[0].Priority() == parserFlt[1].Priority() {
			return nil, fmt.Errorf(
				"found multiple parser with same priority for note type of %s", noteInfo,
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
