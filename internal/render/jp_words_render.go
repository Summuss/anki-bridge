package render

import (
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/model"
)

func init() {
	renderList = append(renderList, jpWordsRender{})
}

type jpWordsRender struct {
}

func (j jpWordsRender) Process(m model.IModel) (*anki.Card, error) {
	jpWord := m.(*model.JPWord)
	return &anki.Card{
		Front: jpWord.Spell,
		Back:  jpWord.Mean,
		Desk:  "test",
	}, nil
}

func (j jpWordsRender) Match(m model.IModel) bool {
	return "jp_words" == m.CollectionName()
}
