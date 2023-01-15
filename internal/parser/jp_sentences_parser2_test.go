package parser

import (
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"reflect"
	"testing"
)

func TestJPSentencesParser2_Check(t *testing.T) {
	type args struct {
		note string
		in1  common.NoteInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				in1: *config.Conf.GetNoteInfoByName(common.NoteType_JPSentences_Name), note: `
- #2 君の一生の思い出、 {{cloze しかと}} {{cloze  見届け}} たぞ。
	- 確と #[[Jp Words]]
		- 〔たしかに〕确实,明确,准确.
		- しかと　2　4
	- 見届ける#[[Jp Words]]
		- 看到，看准，看清；一直看到最后，结束，用眼看，确认。（最後まで見る。また、しっかり見る。）
		- みとどける　0　7
	- this is a cross line explanation
	  xxxx
		- 2 level
		  aaaa
	- this is explanation2
		- xxx
			- 2232
			-

`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				J := JPSentencesParser2{baseParser{}}
				if err := J.Check(tt.args.note, tt.args.in1); (err != nil) != tt.wantErr {
					t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestJPSentencesParser2_Parse(t *testing.T) {
	type args struct {
		note     string
		noteType common.NoteInfo
	}
	word1 := model.JPWord{
		BaseModel:   model.BaseModel{},
		Hiragana:    "できごと",
		Mean:        "(偶发)的事件，变故。（持ち上がった事件・事柄。）",
		Pitch:       "2",
		Spell:       "出来事",
		WordClasses: []int{10, 3},
	}
	word2 := word1
	tests := []struct {
		name    string
		args    args
		want    model.IModel
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				noteType: *config.Conf.GetNoteInfoByName(common.NoteType_JPSentences_Name), note: `
- #2 君の一生の思い出、 {{cloze しかと}} {{cloze  見届け}} たぞ。
	- 出来事#[[Jp Words]]
		- (偶发)的事件，变故。（持ち上がった事件・事柄。）
		- できごと 2 A 3
	- 出来事 #[[Jp Words]]
		- (偶发)的事件，变故。（持ち上がった事件・事柄。）
		- できごと 2 A 3
	- this is a cross line explanation
	  xxxx
		- 2 level
		  aaaa
	- this is explanation2
		- xxx
			- 2232
			-
`,
			},
			want: &model.JPSentence{
				Addition: `- this is a cross line explanation
  xxxx
	- 2 level
	  aaaa
- this is explanation2
	- xxx
		- 2232
		-`,
				Sentence: `君の一生の思い出、 {{cloze しかと}} {{cloze  見届け}} たぞ。`,
				JPWords:  &[]*model.JPWord{&word2, &word1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				J := JPSentencesParser2{
					baseParser{},
				}
				got, err := J.Parse(tt.args.note, tt.args.noteType)
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
