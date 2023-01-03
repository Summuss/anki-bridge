package common

import (
	"fmt"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

func UnIndent(txt string) string {
	r, _ := regexp.Compile(`(?m)^\t`)
	return r.ReplaceAllString(txt, "")
}

func SplitByNoIndentLine(txt string) (*[]string, error) {
	r, _ := regexp.Compile(`(?m)^\S+.*$`)
	split := r.Split(txt, -1)
	if len(strings.TrimSpace(split[0])) > 0 {
		return nil, fmt.Errorf("error near %s", split[0])
	}
	matches := r.FindAllString(txt, -1)
	res := lo.Map(
		matches, func(item string, i int) string {
			return item + split[i+1]
		},
	)
	return &res, nil
}

func PreprocessNote(note string) string {
	replaceMap := map[string]string{
		"＃":  "#",
		"＆":  "&",
		"　":  " ",
		"？":  "?",
		"＠":  "@",
		"１":  "1",
		"２":  "2",
		"３":  "3",
		"４":  "4",
		"５":  "5",
		"６":  "6",
		"７":  "7",
		"８":  "8",
		"９":  "9",
		"‎":  "",
		"\r": "",
	}
	for k, v := range replaceMap {
		note = strings.ReplaceAll(note, k, v)
	}
	return note
}

func SplitWithTrimAndOmitEmpty(s string, step string) []string {
	splits := strings.Split(s, step)
	return lo.Filter(
		splits, func(item string, _ int) bool {
			return len(strings.TrimSpace(item)) > 0
		},
	)

}
