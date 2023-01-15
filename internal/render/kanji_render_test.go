package render

import (
	"github.com/summuss/anki-bridge/internal/model"
	"testing"
)

func Test_renderKanjiBack(t *testing.T) {
	type args struct {
		kanji *model.Kanji
	}
	pron1 := model.KanjiPron{Pron: "イチ", Example: "一度 一座 第一"}
	pron2 := model.KanjiPron{Pron: "イツ", Example: "一般 同一 統一"}
	prons := []*model.KanjiPron{&pron1, &pron2}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{kanji: &model.Kanji{Kanji: "一", Prons: &prons}},
			want: `
	<div id="kanji_41" class="popup_info mincho">
		<table class="char_read">
			<tbody>
<tr>
	<td rowspan="2" class="char"><span>一</span></td>
	<td class="read underline"><span class="pron">イチ</span></td>
	<td class="rei underline">一度 一座 第一</td>
</tr>

<tr>
	<td class="read underline"><span class="pron">イツ</span></td>
	<td class="rei underline">一般 同一 統一</td>
</tr></tbody>
		</table>
	</div>`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := renderKanji(tt.args.kanji); got != tt.want {
					t.Errorf("renderKanji() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
