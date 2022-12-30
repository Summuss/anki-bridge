package util

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

func SplitByNoIndentLine(txt string) ([]string, error) {
	r, _ := regexp.Compile(`(?m)^\S+`)
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
	return res, nil
}
