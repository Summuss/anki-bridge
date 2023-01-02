package render

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/model"
	"regexp"
	"strings"
)

func init() {
	renderList = append(renderList, jpSentencesRender{})
}

type jpSentencesRender struct {
}

func (j jpSentencesRender) Process(m model.IModel) (*anki.Card, error) {
	jpSentence := m.(*model.JPSentence)
	var e error
	words := lo.Map(
		*(jpSentence.JPWords), func(item *model.JPWord, _ int) string {
			card, err := jpWordsRenderIns.Process(item)
			if err != nil {
				e = err
				return ""
			}
			return card.Back
		},
	)
	if e != nil {
		return nil, fmt.Errorf("render word failed,error: %s", e.Error())
	}
	return &anki.Card{
		Front: replaceCloze(jpSentence.Sentence),
		Back:  strings.Join(words, "<hr><br>"),
		Desk:  "test",
	}, nil
}

func replaceCloze(ori string) string {
	r := regexp.MustCompile(`\{\{cloze (.*?)}}`)
	return r.ReplaceAllStringFunc(
		ori, func(s string) string {
			submatches := r.FindStringSubmatch(s)
			return fmt.Sprintf("<span class=\"cloze\">%s</span>", submatches[1])
		},
	)
}

func (j jpSentencesRender) Match(model model.IModel) bool {
	return model.CollectionName() == "jp_sentences"
}
