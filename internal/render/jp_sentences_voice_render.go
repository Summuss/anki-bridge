package render

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"regexp"

	"strings"
)

func init() {
	renderList = append(renderList, jpSentencesVoiceRender{})
}

type jpSentencesVoiceRender struct {
}

func (j jpSentencesVoiceRender) Process(m model.IModel) (*anki.Card, error) {
	jpSentence := m.(*model.JPSentence)
	noteInfo := config.Conf.GetNoteInfoByName(common.NoteType_JPSentences_Name)
	return &anki.Card{
		Front:         fmt.Sprintf("[sound:%s]", (*jpSentence.GetResources())[0].Metadata.FileName),
		Back:          renderBack(jpSentence),
		Desk:          noteInfo.Desk,
		AnkiNoteModel: noteInfo.AnkiNoteModel,
	}, nil
}

func renderBack(m *model.JPSentence) string {
	templStr := `
<div class="jp_sentence">
    <div class="sentence">{{.Sentence}} <br>
        <ol>
			{{.WordMeanings}}
        </ol>
    </div>
</div>
`
	var merr *multierror.Error
	var wordMeanings []string
	if m.JPWords != nil {
		wordMeanings = lo.Map(
			*m.JPWords, func(item *model.JPWord, _ int) string {
				var t = JPWordRenderObj{
					JPWord:    item,
					WordClass: renderWordClasses(item.WordClasses),
					Pitch:     renderPitch(item.Pitch),
				}

				wordTmpl := `
{{- /*gotype: github.com/summuss/anki-bridge/internal/render.JPWordRenderObj*/ -}}
<div class="jp_word">
    <div class="word">
        {{.Spell}} | {{.WordClass}}
    </div>
    <div>
        <span class="hira">{{.Hiragana}}</span> | {{.Pitch }}
    </div>
    <div class="word-meaning">{{.Mean}} </div>
</div>
`
				w, err := renderWordInfoByTempl(wordTmpl, t)
				if err != nil {
					merr = multierror.Append(merr, err)
					return ""
				}
				return fmt.Sprintf("<li>%s</li><br>", w)
			},
		)
	}
	t1 := strings.ReplaceAll(templStr, "{{.Sentence}}", replaceHighlight(m.Sentence))
	return strings.ReplaceAll(t1, "{{.WordMeanings}}", strings.Join(wordMeanings, "\n"))
}

func (j jpSentencesVoiceRender) Match(model model.IModel) bool {
	return model.GetNoteTypeName() == common.NoteType_JPSentences_Voice_Name
}

func replaceHighlight(ori string) string {
	r := regexp.MustCompile(`\[\[#(\S+)]]==(.*?)==`)
	return r.ReplaceAllStringFunc(
		ori, func(s string) string {
			submatches := r.FindStringSubmatch(s)
			return fmt.Sprintf("<span class=\"%s\">%s</span>", submatches[1], submatches[2])
		},
	)
}
