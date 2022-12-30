package parser

import (
	"fmt"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/util"
	"golang.org/x/exp/slices"
	"regexp"
	"strconv"
)

var jpWordsParser = JPWordsParser{}

type JPWordsParser struct {
}

func (w JPWordsParser) NoteName() string {
	return "Jp Words"
}

func (w JPWordsParser) Split(rowNotes string) ([]string, error) {
	return util.SplitByNoIndentLine(rowNotes)
}

func (w JPWordsParser) Check(note string) error {
	notePreproc := util.PreprocessNote(note)
	r, _ := regexp.Compile(`(?m)\A\s*^-\s*(?P<word>\S+)$\n^\t-\s*(?P<meaning>\S+.*)$\n^\t-\s*(?P<hiragana>\S+)\s+(?P<pitch>\d)\s+(?P<classes>.+)$\s*\z`)
	if !r.MatchString(notePreproc) {
		return fmt.Errorf("synatx error in word\n%s", note)
	}
	submatches := r.FindStringSubmatch(notePreproc)
	classesStr := submatches[r.SubexpIndex("classes")]
	classes := util.SplitWithTrimAndOmitEmpty(classesStr, " ")
	if len(classes) == 0 {
		return fmt.Errorf("word classes not found in word\n%s", note)
	}
	for _, class := range classes {
		err := checkWordClass(class)
		if err != nil {
			return fmt.Errorf("%s in word\n%s", err.Error(), note)
		}
	}
	return nil
}

func (w JPWordsParser) Parse(note string) (model.Model, error) {
	//TODO implement me
	panic("implement me")
}
func checkWordClass(class string) error {
	n, err := strconv.ParseInt(class, 16, 32)
	if err != nil {
		return fmt.Errorf("incorrect class: %s", class)
	}
	classSet := []int{2, 3, 1, 5, 4, 6, 7, 12, 11, 10, 8, 9}
	if !slices.Contains(classSet, int(n)) {
		return fmt.Errorf("incorrect class: %s", class)
	}
	return nil
}
