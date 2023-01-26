package parser

import (
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"reflect"
	"testing"
)

func TestJPSentencesVoiceParser_Parse(t *testing.T) {
	type args struct {
		note     string
		noteType *common.NoteInfo
	}
	tests := []struct {
		name           string
		args           args
		want           model.IModel
		wantMiddleInfo interface{}
		wantErr        bool
	}{
		{
			name: "1",
			args: args{
				note: `
- #FILENAME [2023-01-16][20-52-30].mp3 #FILENAME [2023-01-16][20-42-56].mp3
	- あれが[[#red]]==偏屈==じゃなくて何なんだ！
	- 偏屈#[[Jp Words]]
		- 怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）
		- へんくつ 1 1
	- 偏屈#[[Jp Words]]
		- 怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）
		- へんくつ 1 1
`,
			},
			want: &model.JPSentence{
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
			},
			wantMiddleInfo: &jPSentencesVoiceMiddleInfo{
				filePaths: []string{
					"D:\\Documents\\voice-records/processed/[2023-01-16][20-52-30].mp3",
					"D:\\Documents\\voice-records/[2023-01-16][20-42-56].mp3",
				},
			},
			wantErr: false,
		},
		{
			name: "only sentence",
			args: args{
				note: `
- #FILENAME [2023-01-16][20-52-30].mp3
	- あれが[[#red]]==偏屈==じゃなくて何なんだ！
`,
			},
			want: &model.JPSentence{
				Sentence: `あれが[[#red]]==偏屈==じゃなくて何なんだ！`,
				JPWords:  &[]*model.JPWord{},
			},
			wantErr: false,
		},
		{
			name: "only words",
			args: args{
				note: `
- #FILENAME [2023-01-16][20-52-30].mp3
	- 偏屈#[[Jp Words]]
		- 怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）
		- へんくつ 1 1
	- 偏屈#[[Jp Words]]
		- 怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）
		- へんくつ 1 1
`,
			},
			want: &model.JPSentence{
				Sentence: ``,
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
			},
			wantErr: false,
		},
		{
			name: "word check error",
			args: args{
				note: `
- #FILENAME [2023-01-16][20-52-30].mp3
	- あれが[[#red]]==偏屈==じゃなくて何なんだ！
	- 偏屈#[[Jp Words]]
		- 怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）
		- へんくつ 1 x
	- 偏屈#[[Jp Words]]
		- 怪癖，乖僻，顽固，别扭；古怪；孤僻。（性質が、ねじけていること）
		- へんくつ 1 1
`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only file",
			args: args{
				note: `
- #FILENAME [2023-01-25][01-12-05].mp3 #FILENAME [2023-01-25][01-13-09].mp3
`,
			},
			want: &model.JPSentence{
				Sentence: ``,
				JPWords:  &[]*model.JPWord{},
			},
			wantErr: false,
			wantMiddleInfo: &jPSentencesVoiceMiddleInfo{
				filePaths: []string{
					"D:\\Documents\\voice-records/[2023-01-25][01-12-05].mp3",
					"D:\\Documents\\voice-records/[2023-01-25][01-13-09].mp3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				J := JPSentencesVoiceParser{}
				got, err := J.MiddleParse(tt.args.note, tt.args.noteType)
				words := got.(*model.JPSentence).JPWords
				for _, word := range *words {
					word.SetNoteInfo(nil)
					word.SetParser(nil)
				}
				if (err != nil) != tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				tt.want.SetMiddleInfo(tt.wantMiddleInfo)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parse() got = %v", got)
					t.Errorf("Parse() want %v", tt.want)

				}
			},
		)
	}
}

func TestJPSentencesVoiceParser_PostParse(t *testing.T) {
	type args struct {
		iModel model.IModel
	}
	jpSentence := &model.JPSentence{
		Sentence: ``,
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
	jpSentence.SetMiddleInfo(
		&jPSentencesVoiceMiddleInfo{
			filePaths: []string{
				"D:\\Documents\\voice-records/processed/[2023-01-16][20-52-30].mp3",
				"D:\\Documents\\voice-records/[2023-01-16][20-42-56].mp3",
			},
		},
	)
	jpSentence.SetNoteInfo(config.Conf.GetNoteInfoByName(common.NoteType_JPSentences_Voice_Name))
	tests := []struct {
		name    string
		args    args
		want    model.IModel
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{jpSentence},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				J := JPSentencesVoiceParser{}
				got, err := J.PostParse(tt.args.iModel)
				if (err != nil) != tt.wantErr {
					t.Errorf("PostParse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(*got.GetResources()) != 2 {
					t.Errorf("resource not filled")
				}
			},
		)
	}
}

func Test_moveVoiceFile(t *testing.T) {
	err := moveVoiceFile("D:\\Documents\\voice-records\\新規 テキスト ドキュメント.txt")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = moveVoiceFile("D:\\Documents\\voice-records\\processed\\新規 テキスト ドキュメント.txt")
	if err != nil {
		t.Errorf(err.Error())
	}

}
