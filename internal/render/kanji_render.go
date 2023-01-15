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
	noteInfo := config.Conf.GetNoteInfoByName(common.NoteType_Kanji_Name)
	return &anki.Card{
		Front:         renderKanji(kanji),
		Back:          "",
		Desk:          noteInfo.Desk,
		AnkiNoteModel: noteInfo.AnkiNoteModel,
	}, nil
}

func renderKanji(kanji *model.Kanji) string {
	size := len(*kanji.Prons)
	trs := lo.Map(
		*kanji.Prons, func(item *model.KanjiPron, i int) string {
			if i == 0 {
				tr := fmt.Sprintf(
					`
<tr>
	<td rowspan="%d" class="char"><span>%s</span></td>
	<td class="read underline"><span class="pron">%s</span></td>
	<td class="rei underline">%s</td>
</tr>`,
					size, kanji.Kanji, item.Pron, item.Example,
				)
				return tr

			}
			tr := fmt.Sprintf(
				`
<tr>
	<td class="read underline"><span class="pron">%s</span></td>
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
