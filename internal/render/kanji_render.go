package render

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"strings"
)

func init() {
	renderList = append(renderList, kanjiRender{})

}

type kanjiRender struct {
}

func (j kanjiRender) Process(m model.IModel) (*anki.Card, error) {
	kanji := m.(*model.Kanji)
	return &anki.Card{
		Front: kanji.Kanji,
		Back:  renderKanjiBack(kanji),
		Desk:  config.Conf.NoteType2Desk[common.NoteType_Kanji],
	}, nil
}

func renderKanjiBack(kanji *model.Kanji) string {
	trs := lo.Map(
		*kanji.Prons, func(item *model.KanjiPron, i int) string {
			if i == 0 {
				tr := fmt.Sprintf(
					`
<tr>
	<td rowspan="5" class="char"><span>%s</span></td>
	<td class="read underline">%s</td>
	<td class="rei underline">%s</td>
</tr>`,
					kanji.Kanji, item.Pron, item.Example,
				)
				return tr

			}
			tr := fmt.Sprintf(
				`
<tr>
	<td class="read underline">%s</td>
	<td class="rei underline">%s</td>
</tr>`,
				item.Pron, item.Example,
			)
			return tr
		},
	)
	return fmt.Sprintf(
		`
	<div id="kanji_41" class="popup_info mincho">
		<table class="char_read">
			<tbody>%s</tbody>
		</table>
	</div>`, strings.Join(trs, "\n"),
	)
}

func (j kanjiRender) Match(m model.IModel) bool {
	return m.CollectionName() == model.KanjiCollectionName
}
