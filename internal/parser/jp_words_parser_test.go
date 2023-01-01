package parser

import (
	"github.com/summuss/anki-bridge/internal/model"
	"reflect"
	"testing"
)

func TestJPWordsParser_Split(t *testing.T) {
	type args struct {
		rowNotes string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "1", args: args{
				rowNotes: `
W1
	1-1
	1-2
W2
	2-1
`,
			}, want: []string{
				`W1
	1-1
	1-2
`, `W2
	2-1
`,
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				w := JPWordsParser{}
				got, err := w.Split(tt.args.rowNotes)
				if (err != nil) != tt.wantErr {
					t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Split() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestJPWordsParser_Check(t *testing.T) {
	type args struct {
		note string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				note: `- 出来事
	- (偶发)的事件，变故。（持ち上がった事件・事柄。）
	- できごと 2 １ 3`,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				note: `- 出来事
	- (偶发)的事件，变故。（持ち上がった事件・事柄。）
	- できごと 2 １ `,
			},
			wantErr: false,
		},

		{
			name: "3",
			args: args{
				note: `- 出来事
	- (偶发)的事件，变故。（持ち上がった事件・事柄。）
	- できごと 2 １ X`,
			},
			wantErr: true,
		},
		{
			name: "start and end with space",
			args: args{
				note: `
- 出来事
	- (偶发)的事件，变故。（持ち上がった事件・事柄。）
	- できごと 2 １ 
`,
			},
			wantErr: false,
		},
		{
			name: "meaning contain space",
			args: args{
				note: `
- 出来事
	- (偶发)的事件，  变故。（持ち上がった事件・事柄。）
	- できごと 2 １ 
`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				w := JPWordsParser{}
				if err := w.Check(tt.args.note); (err != nil) != tt.wantErr {
					t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestJPWordsParser_Parse(t *testing.T) {
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
				AnkiNoteId:  0,
				Hiragana:    "できごと",
				Mean:        "(偶发)的事件，变故。（持ち上がった事件・事柄。）",
				Pitch:       "2",
				Spell:       "出来事",
				WordClasses: []int{10, 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				w := JPWordsParser{}
				got, err := w.Parse(tt.args.note)
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
