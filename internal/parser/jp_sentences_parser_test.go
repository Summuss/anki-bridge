package parser

import (
	"github.com/summuss/anki-bridge/internal/model"
	"reflect"
	"testing"
)

func TestJPSentencesParser_Check(t *testing.T) {
	type args struct {
		note string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "start and end with non-space",
			args: args{
				note: `- 原因不明の病に {{cloze 冒される}} ので
	- 冒す
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
	- 冒す2
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
	- 冒す3
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7`,
			},
			wantErr: false,
		},
		{
			name: "start and end with space",
			args: args{
				note: `
- 原因不明の病に {{cloze 冒される}} ので
	- 冒す
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
	- 冒す2
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
	- 冒す3
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
`,
			},
			wantErr: false,
		},
		{
			name: "syntax err",
			args: args{
				note: `
- 原因不明の病に {{cloze 冒される}} ので
	- 冒す
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
	- 冒す2
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
	- 冒す3
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
xx
`,
			},
			wantErr: true,
		},
		{
			name: "word error",
			args: args{
				note: `
- 原因不明の病に {{cloze 冒される}} ので
	- 冒す
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
	- 冒す2
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　X
	- 冒す3
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7
xx
`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				J := JPSentencesParser{}
				if err := J.Check(tt.args.note); (err != nil) != tt.wantErr {
					t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestJPSentencesParser_Parse(t *testing.T) {
	type args struct {
		note string
	}
	word := &model.JPWord{
		BaseModel:   model.BaseModel{},
		Hiragana:    "おかす",
		Mean:        "侵袭；患（病）。（害を及ぼす。）",
		Pitch:       "2",
		Spell:       "冒す",
		WordClasses: []int{7},
	}
	tests := []struct {
		name    string
		args    args
		want    model.IModel
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				note: `
- 原因不明の病に {{cloze 冒される}} ので
	- 冒す
		- 侵袭；患（病）。（害を及ぼす。）
		- おかす　2　7

`,
			},
			want: &model.JPSentence{
				BaseModel:   model.BaseModel{},
				Sentence:    "原因不明の病に {{cloze 冒される}} ので",
				Explanation: "",
				JPWords:     &[]*model.JPWord{word},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				J := JPSentencesParser{}
				got, err := J.Parse(tt.args.note)
				if (err != nil) != tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parse() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
