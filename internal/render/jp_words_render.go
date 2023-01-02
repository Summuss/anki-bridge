package render

import (
	"bytes"
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/model"
	"html/template"
	"io"
	"strings"
)

var jpWordsRenderIns = jpWordsRender{}

func init() {
	renderList = append(renderList, jpWordsRenderIns)
}

type jpWordsRender struct {
}

func (j jpWordsRender) Process(m model.IModel) (*anki.Card, error) {
	jpWord := m.(*model.JPWord)
	templStr := `
<div class="jp_word">
    <div class="sentence">
        {{.Spell}} | {{.WordClass}}
    </div>
    <div>
        <span class="hira">{{.Hiragana}}</span> | {{.Pitch }} {{.Sound}}
    </div>
</div>
`
	var t = struct {
		*model.JPWord
		WordClass string
		Pitch     string
		Sound     string
	}{
		JPWord:    jpWord,
		WordClass: renderWordClasses(jpWord.WordClasses),
		Pitch:     renderPitch(jpWord.Pitch),
		Sound:     renderSounds(jpWord.GetResources()),
	}
	temp, err := template.New("JPWord").Parse(templStr)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = temp.Execute(&buf, t)
	if err != nil {
		return nil, err
	}
	bts, err := io.ReadAll(&buf)
	if err != nil {
		return nil, err
	}

	return &anki.Card{
		Front: jpWord.Mean,
		Back:  string(bts),
		Desk:  "Japanese::Words",
	}, nil
}

func renderPitch(pitch string) string {
	if pitch == "0" {
		return "◎"
	} else if pitch == "1" {
		return "①"
	} else if pitch == "2" {
		return "②"
	} else if pitch == "3" {
		return "③"
	} else if pitch == "4" {
		return "④"
	} else if pitch == "5" {
		return "⑤"
	} else if pitch == "6" {
		return "⑥"
	} else if pitch == "7" {
		return "⑦"
	} else if pitch == "8" {
		return "⑧"
	} else if pitch == "9" {
		return "⑨"
	} else {
		return "?"
	}
}

func renderWordClasses(wcs []int) string {
	symbol2desc := map[int]string{
		2:  "形容詞",
		3:  "形容動詞",
		1:  "名",
		5:  "動",
		4:  "副",
		6:  "自動",
		7:  "他動",
		12: "サ変",
		11: "五段活用",
		10: "一段活用",
		8:  "接続",
		9:  "接尾",
	}
	res := lo.Map(
		wcs, func(item int, _ int) string {
			return symbol2desc[item]
		},
	)
	return strings.Join(res, "・")
}
func renderSounds(rs *[]model.Resource) string {
	var t []string
	size := len(*rs)
	for i := 0; i < size; i++ {
		fileName := (*rs)[i].Metadata.FileName
		t = append(t, fmt.Sprintf("[sound:%s]", fileName))
	}
	return strings.Join(t, " ")
}

func (j jpWordsRender) Match(m model.IModel) bool {
	return "Jp Words" == m.GetNoteType()
}
