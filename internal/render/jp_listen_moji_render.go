package render

import (
	"fmt"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/dto"
	"github.com/summuss/anki-bridge/internal/model"
	"strings"
)

func init() {
	renderList = append(renderList, jpListenMojiRender{})

}

type jpListenMojiRender struct {
}

func (j jpListenMojiRender) Process(m model.IModel) (*anki.Card, error) {
	jpListenMoji := m.(*model.JpListen[dto.MojiBookItem])
	noteInfo := config.Conf.GetNoteInfoByName(common.NoteType_JpListenMoji_Name)
	templStr := `
<div class="jp_word">
    <div class="sentence">
        {{.Spell}} 
    </div>
    <div>
        <span class="hira">{{.Hiragana}}</span> | {{.Pitch}} {{.Sound}}
    </div>
    <div>
        <span>{{.Excerpt}} </span>
    </div>
</div>
`
	sounds := renderSounds(jpListenMoji.GetResources())
	back := strings.Replace(templStr, "{{.Spell}}", jpListenMoji.ExtInfo.Spell, 1)
	back = strings.Replace(back, "{{.Hiragana}}", jpListenMoji.ExtInfo.Pron, 1)
	back = strings.Replace(back, "{{.Pitch}}", jpListenMoji.ExtInfo.Accent, 1)
	back = strings.Replace(back, "{{.Sound}}", sounds, -1)
	back = strings.Replace(back, "{{.Excerpt}}", jpListenMoji.ExtInfo.Excerpt, 1)

	return &anki.Card{
		Front:         fmt.Sprintf("%s %s %s", sounds, sounds, sounds),
		Back:          back,
		Desk:          noteInfo.Desk,
		AnkiNoteModel: noteInfo.AnkiNoteModel,
	}, nil

}

func (j jpListenMojiRender) Match(m model.IModel) bool {
	return m.GetNoteTypeName() == common.NoteType_JpListenMoji_Name
}
