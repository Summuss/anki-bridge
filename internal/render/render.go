package render

import (
	"fmt"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"reflect"
	"strings"
)

var renderList []iRender

type iRender interface {
	Process(model.IModel) (*anki.Card, error)
	Match(model.IModel) bool
}

func Render(m model.IModel) (*anki.Card, error) {
	for _, r := range renderList {
		if r.Match(m) {
			card, err := r.Process(m)
			if card != nil && !config.Conf.RealMode {
				card.Desk = "test"
			}
			return card, err
		}
	}
	return nil, fmt.Errorf(
		"can't find render for model type %s", reflect.ValueOf(m).Type().Elem().Name(),
	)

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
