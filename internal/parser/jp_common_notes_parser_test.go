package parser

import (
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/model"
	"reflect"
	"testing"
)

func TestJPCommonNoteParser_MiddleParse(t *testing.T) {
	type args struct {
		note     string
		noteType *common.NoteInfo
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
				`- ～だろうが
	- 子供だろうが、大人だろうが、法を守らなければならないのは同じだ。
	- 生きようが死のうが……
	- 寒かろうが、暑かろうが……
`, nil,
			},
			want: &model.JPCommonNote{
				Front: `～だろうが`, Back: `- 子供だろうが、大人だろうが、法を守らなければならないのは同じだ。
- 生きようが死のうが……
- 寒かろうが、暑かろうが……`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				J := jpCommonNoteParser
				got, err := J.MiddleParse(tt.args.note, tt.args.noteType)
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
