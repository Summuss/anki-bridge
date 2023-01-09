package render

import (
	"bytes"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"html/template"
	"io"
)

func init() {
	renderList = append(renderList, jpRecognitionRender{})

}

type jpRecognitionRender struct {
}

func (j jpRecognitionRender) Process(m model.IModel) (*anki.Card, error) {
	jpWord := m.(*model.JPWord)
	templStr := `
<div class="jp_recognition">
    <div class="sentence jp_word">
        {{.Spell}} | {{.WordClass}}
    </div>
    <div>
        <span class="hira">{{.Hiragana}}</span> | {{.Pitch }} {{.Sound}}
        <div><span class="meaning"> {{.Mean}} </span></div>
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
		Front: jpWord.Spell,
		Back:  string(bts),
		Desk:  config.Conf.GetNoteInfoByName(common.NoteType_JPRecognition_Name).Desk,
	}, nil
}

func (j jpRecognitionRender) Match(m model.IModel) bool {
	return m.GetNoteTypeName() == common.NoteType_JPRecognition_Name
}
