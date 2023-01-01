package parser

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/model"
	"github.com/summuss/anki-bridge/internal/util"
	"golang.org/x/exp/slices"
	"regexp"
	"strconv"
	"strings"
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

func (w JPWordsParser) Parse(note string) (model.IModel, error) {
	notePreproc := util.PreprocessNote(note)
	r, _ := regexp.Compile(`(?m)\A\s*^-\s*(?P<word>\S+)$\n^\t-\s*(?P<meaning>\S+.*)$\n^\t-\s*(?P<hiragana>\S+)\s+(?P<pitch>\d)\s+(?P<classes>.+)$\s*\z`)
	submatches := r.FindStringSubmatch(notePreproc)
	classesStr := submatches[r.SubexpIndex("classes")]
	classes := util.SplitWithTrimAndOmitEmpty(classesStr, " ")
	word := submatches[r.SubexpIndex("word")]
	meaning := submatches[r.SubexpIndex("meaning")]
	hiragana := submatches[r.SubexpIndex("hiragana")]
	pitch := submatches[r.SubexpIndex("pitch")]
	classes_int := lo.Map(
		classes, func(item string, _ int) int {
			res, _ := strconv.ParseInt(item, 16, 32)
			return int(res)
		},
	)
	maleURL, femaleURL, err := getTTSURL(word)
	if err != nil {
		return nil, err
	}
	data1, err := util.CurlGetData(maleURL)
	if err != nil {
		return nil, fmt.Errorf("download tts from %s failed,error:\n%s", maleURL, err.Error())
	}
	data2, err := util.CurlGetData(femaleURL)
	if err != nil {
		return nil, fmt.Errorf("download tts from %s failed,error:\n%s", femaleURL, err.Error())
	}

	resource1 := &model.Resource{
		Metadata: model.ResourceMetadata{
			FileName: word + "-male.mp3", ResourceType: model.Sound, ExtName: ".mp3",
		},
	}
	resource1.SetData(*data1)

	resource2 := &model.Resource{
		Metadata: model.ResourceMetadata{
			FileName: word + "-female.mp3", ResourceType: model.Sound, ExtName: ".mp3",
		},
	}
	resource2.SetData(*data2)

	jpWord := &model.JPWord{
		Hiragana: hiragana, Mean: meaning, Pitch: pitch, Spell: word, WordClasses: classes_int,
	}
	jpWord.SetResources(&[]model.Resource{*resource1, *resource2})
	return jpWord, nil

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

func getTTSURL(txt string) (maleURL string, femaleURL string, err error) {
	txt = "\"" + txt + "\""
	//FIXME
	jsFile := "/home/summus/Code/script/python-script/jp_study/oddcast_api.js"
	res1, err := util.Exec("/usr/bin/node", jsFile, txt, "takeru")

	if err != nil {
		return "", "", fmt.Errorf(
			"\"node %s %s %s\" exec failed,%s", jsFile, txt, "takeru", err.Error(),
		)
	}
	res2, err := util.Exec("/usr/bin/node", jsFile, txt, "sayaka")
	if err != nil {
		return "", "", fmt.Errorf(
			"\"node %s %s %s\" exec failed,%s", jsFile, txt, "sayaka", err.Error(),
		)
	}
	res1 = strings.ReplaceAll(res1, "\n", "")
	res2 = strings.ReplaceAll(res2, "\n", "")
	return res1, res2, nil

}
