package render

import (
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/model"
	"reflect"
	"testing"
)

func Test_jpSentencesRender_Process(t *testing.T) {
	type args struct {
		m model.IModel
	}

	r1 := model.Resource{}
	r2 := model.Resource{}
	r1.Metadata.FileName = "衰える-male.mp3"
	r2.Metadata.FileName = "衰える-female.mp3"
	rs := []model.Resource{r1, r2}
	baseModel := model.BaseModel{}
	baseModel.SetResources(&rs)

	word1 := model.JPWord{
		BaseModel:   baseModel,
		Hiragana:    "あまつさえ",
		Mean:        "而且并且",
		Pitch:       "1",
		Spell:       "衰える",
		WordClasses: []int{3, 12},
	}
	word2 := word1
	words := []*model.JPWord{&word1, &word2}

	tests := []struct {
		name    string
		args    args
		want    *anki.Card
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				m: &model.JPSentence{
					Sentence: "君の一生の思い出、 {{cloze しかと}} {{cloze  見届け}} たぞ。",
					JPWords:  &words,
					Addition: `- this is a cross line explanation
  xxxx
	- 2 level
	  aaaa
- this is explanation2
	- xxx
		- 2232
		-`,
				},
			},
			want: &anki.Card{
				Front: `君の一生の思い出、 <span class="cloze">しかと</span> <span class="cloze"> 見届け</span> たぞ。`,
				Back: `
<div class="jp_word">
    <div class="sentence">
        衰える | 形容動詞・サ変
    </div>
    <div>
        <span class="hira">あまつさえ</span> | ① [sound:衰える-male.mp3] [sound:衰える-female.mp3]
    </div>
</div>
<hr><br>
<div class="jp_word">
    <div class="sentence">
        衰える | 形容動詞・サ変
    </div>
    <div>
        <span class="hira">あまつさえ</span> | ① [sound:衰える-male.mp3] [sound:衰える-female.mp3]
    </div>
</div>
<div class="addition"><ul>
<li>this is a cross line explanation
xxxx

<ul>
<li>2 level
aaaa</li>
</ul></li>
<li>this is explanation2

<ul>
<li>xxx

<ul>
<li>2232
-</li>
</ul></li>
</ul></li>
</ul>
</div>`,
				Desk: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				j := jpSentencesRender{}
				got, err := j.Process(tt.args.m)
				if (err != nil) != tt.wantErr {
					t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Process() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_renderFront(t *testing.T) {
	type args struct {
		m *model.JPSentence
	}
	word1 := model.JPWord{
		Hiragana:    "あまつさえ",
		Mean:        "而且并且",
		Pitch:       "1",
		Spell:       "衰える",
		WordClasses: []int{3, 12},
	}

	words := []*model.JPWord{&word1}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				m: &model.JPSentence{
					Sentence: "君の一生の思い出、 {{cloze しかと}} {{cloze 見届け}} たぞ。",
					JPWords:  &words,
				},
			},
			want: `
<div class="jp_sentence">
    <div class="sentence">君の一生の思い出、<span class="cloze">しかと</span><span class="cloze">見届け</span>たぞ。 <br>
        <ol>
			<li>而且并且</li>
        </ol>
    </div>
</div>
`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := renderFront(tt.args.m); got != tt.want {
					t.Errorf("renderFront() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
