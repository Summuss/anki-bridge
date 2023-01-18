package render

import (
	"github.com/summuss/anki-bridge/internal/model"
	"testing"
)

func Test_renderBack(t *testing.T) {
	type args struct {
		m *model.JPSentence
	}
	jpSentence := &model.JPSentence{
		Sentence: `あれが[[#red]]==偏屈==じゃなくて何なんだ！`,
		JPWords: &[]*model.JPWord{
			{
				Hiragana:    "へんくつ",
				Mean:        "怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）",
				Pitch:       "1",
				Spell:       "偏屈",
				WordClasses: []int{1},
			},
			{
				Hiragana:    "へんくつ",
				Mean:        "怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）",
				Pitch:       "1",
				Spell:       "偏屈",
				WordClasses: []int{1},
			},
		},
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				jpSentence,
			},
			want: `
<div class="jp_sentence">
    <div class="sentence">あれが<span class="red">偏屈</span>じゃなくて何なんだ！ <br>
        <ol>
			<li><div class="jp_word">
    <div class="word">
        偏屈 | 名
    </div>
    <div>
        <span class="hira">へんくつ</span> | ①
    </div>
    <div class="word-meaning">怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること） </div>
</div>
</li><br>
<li><div class="jp_word">
    <div class="word">
        偏屈 | 名
    </div>
    <div>
        <span class="hira">へんくつ</span> | ①
    </div>
    <div class="word-meaning">怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること） </div>
</div>
</li><br>
        </ol>
    </div>
</div>
`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := renderBack(tt.args.m); got != tt.want {
					t.Errorf("renderBack() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
