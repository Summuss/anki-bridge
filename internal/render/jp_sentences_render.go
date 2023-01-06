package render

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
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
	var addiHTML string
	if len(jpSentence.Addition) > 0 {
		bs := markdown.ToHTML([]byte(jpSentence.Addition), nil, nil)
		addiHTML = string(bs)
		addiHTML = fmt.Sprintf(`<div class="addition">%s</div>`, addiHTML)
	}
	return &anki.Card{
		Front: renderFront(jpSentence),
		Back:  strings.Join(words, "<hr><br>") + addiHTML,
		Desk:  config.Conf.NoteType2Desk[common.NoteType_JPSentences],
	}, nil
}

func renderFront(m *model.JPSentence) string {
	templStr := `
<div class="jp_sentence">
    <div class="sentence">{{.Sentence}} <br>
        <ol>
			{{.WordMeanings}}
        </ol>
    </div>
</div>
`
	wordMeanings := lo.Map(
		*m.JPWords, func(item *model.JPWord, _ int) string {
			return fmt.Sprintf("<li>%s</li>", item.Mean)
		},
	)

	t1 := strings.ReplaceAll(templStr, "{{.Sentence}}", replaceCloze(m.Sentence))
	return strings.ReplaceAll(t1, "{{.WordMeanings}}", strings.Join(wordMeanings, "\n"))
}

func replaceCloze(ori string) string {
	r := regexp.MustCompile(`\s?\{\{cloze (.*?)}}\s?`)
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
