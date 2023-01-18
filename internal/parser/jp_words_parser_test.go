package parser

import (
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
	"reflect"
	"strings"
	"testing"
)

func TestJPWordsParser_MiddleParse(t *testing.T) {
	type args struct {
		note string
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
- 出来事
	- (偶发)的事件，变故。（持ち上がった事件・事柄。）
	- できごと 2 A 3

`,
			},
			want: &model.JPWord{
				BaseModel:   model.BaseModel{},
				Hiragana:    "できごと",
				Mean:        "(偶发)的事件，变故。（持ち上がった事件・事柄。）",
				Pitch:       "2",
				Spell:       "出来事",
				WordClasses: []int{10, 3},
			},
			wantErr: false,
		},
		{
			name: "word class error",
			args: args{
				note: `
- 出来事
	- (偶发)的事件，变故。（持ち上がった事件・事柄。）
	- できごと 2 A X

`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				w := JPWordsParser{}
				got, err := w.MiddleParse(
					tt.args.note, config.Conf.GetNoteInfoByName(common.NoteType_JPWords_Name),
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("MiddleParse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MiddleParse() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_getTTSURL(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name          string
		args          args
		wantMaleURL   string
		wantFemaleURL string
		wantErr       bool
	}{
		{
			name:          "1",
			args:          args{txt: `日本語`},
			wantMaleURL:   "",
			wantFemaleURL: "",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotMaleURL, gotFemaleURL, err := getTTSURL(tt.args.txt)
				if (err != nil) != tt.wantErr {
					t.Errorf("getTTSURL() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.HasPrefix(gotMaleURL, "https") {
					t.Errorf("getTTSURL() gotMaleURL = %v, want %v", gotMaleURL, tt.wantMaleURL)
				}
				if !strings.HasPrefix(gotFemaleURL, "https") {
					t.Errorf(
						"getTTSURL() gotFemaleURL = %v, want %v", gotFemaleURL, tt.wantFemaleURL,
					)
				}
			},
		)
	}
}
