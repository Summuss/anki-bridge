package render

import (
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/model"
)

func init() {
	//renderList = append(renderList, demosRender{})

}

type demosRender struct {
}

func (j demosRender) Process(model model.IModel) (*anki.Card, error) {
	panic("implement me")
}

func (j demosRender) Match(model model.IModel) bool {
	panic("implement me")
}