package render

import (
	"github.com/gomarkdown/markdown"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
)

func init() {
	renderList = append(renderList, JPCommonNoteRender{})

}

type JPCommonNoteRender struct {
}

func (j JPCommonNoteRender) Process(m model.IModel) (*anki.Card, error) {
	jpCommonNote := m.(*model.JPCommonNote)
	noteInfo := config.Conf.GetNoteInfoByName(common.NoteType_JPCommonNotes_Name)
	bs := markdown.ToHTML([]byte(jpCommonNote.Back), nil, nil)

	return &anki.Card{
		Front:         jpCommonNote.Front,
		Back:          string(bs),
		Desk:          noteInfo.Desk,
		AnkiNoteModel: noteInfo.AnkiNoteModel,
	}, nil

}

func (j JPCommonNoteRender) Match(m model.IModel) bool {
	return m.GetNoteTypeName() == common.NoteType_JPCommonNotes_Name
}
