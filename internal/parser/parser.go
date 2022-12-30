package parser

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/util"
	"regexp"
	"strings"
	"sync"
)

type Parser interface {
	NoteName() string
	Split(string) ([]string, error)
	Check(string) error
	Parse(string) (model.Model, error)
}

var (
	parsers = &[]Parser{jpWordsParser, jpSentencesParser}
)

func splitByNoteType(content string) (map[string]string, error) {
	content = content + "\n"
	r2, _ := regexp.Compile(`(?m)^\S+?$\n`)
	splits := r2.Split(content, -1)
	if len(strings.TrimSpace(splits[0])) > 0 {
		return nil, fmt.Errorf("found error near %s", splits[0])
	}
	matches := r2.FindAllString(content, -1)

	noteTypeName2SubTxt := make(map[string]string)
	for i, match := range matches {
		noteTypeName := strings.TrimSpace(match)
		note := strings.TrimSpace(util.UnIndent(splits[i+1]))

		if len(noteTypeName) == 0 || len(note) == 0 {
			continue
		}
		noteTypeName2SubTxt[noteTypeName] += "\n" + note
	}
	return noteTypeName2SubTxt, nil
}

func CheckInput(txt string) error {
	noteTypeName2SubTxt, err := splitByNoteType(txt)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	size := len(noteTypeName2SubTxt)
	if size == 0 {
		return nil
	}
	wg.Add(size)

	errList := util.SafeList[error]{}
	for noteName, subTxt := range noteTypeName2SubTxt {
		noteName := noteName
		subTxt := subTxt
		go func() {
			defer wg.Done()
			parser, err := findParser(noteName)
			if err != nil {
				errList.Add(err)
				return
			}
			notes, err := parser.Split(subTxt)
			if err != nil {
				errList.Add(err)
				return
			}
			err = checkNotes(notes, parser)
			if err != nil {
				errList.Add(err)
				return
			}

		}()
	}
	wg.Wait()
	return util.MergeErrors(errList.ToSlice())
}

func checkNotes(notes []string, p Parser) error {
	if len(notes) == 0 {
		return nil

	}
	var wg sync.WaitGroup
	wg.Add(len(notes))
	errList := util.SafeList[error]{}

	for _, note := range notes {
		note := note
		go func() {
			defer wg.Done()
			err := p.Check(note)
			if err != nil {
				errList.Add(err)
			}
		}()
	}
	wg.Wait()
	return util.MergeErrors(errList.ToSlice())
}

func findParser(noteName string) (Parser, error) {
	parserFlt := lo.Filter(
		*parsers, func(item Parser, index int) bool {
			return item.NoteName() == noteName
		},
	)
	if len(parserFlt) < 1 {
		return nil, fmt.Errorf("can't found parser for note type of %s", noteName)
	} else if len(parserFlt) > 1 {
		return nil, fmt.Errorf("found %d parser for note type of %s", len(parserFlt), noteName)
	}
	p := parserFlt[0]
	return p, nil
}
