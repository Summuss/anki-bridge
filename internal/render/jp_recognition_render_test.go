package render

import (
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/model"
	"reflect"
	"testing"
)

func Test_jpRecognitionRender_Process(t *testing.T) {
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

	tests := []struct {
		name    string
		args    args
		want    *anki.Card
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				m: &model.JPWord{
					BaseModel:   baseModel,
					Hiragana:    "あまつさえ",
					Mean:        "而且并且",
					Pitch:       "1",
					Spell:       "衰える",
					WordClasses: []int{3, 12},
				},
			},
			want: &anki.Card{
				Front: "衰える",
				Back: `
<div class="jp_recognition">
    <div class="sentence jp_word">
        衰える | 形容動詞・サ変
    </div>
    <div>
        <span class="hira">あまつさえ</span> | ① [sound:衰える-male.mp3] [sound:衰える-female.mp3]
        <div><span class="meaning"> 而且并且 </span></div>
    </div>
</div>
`,
				Desk: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				j := jpRecognitionRender{}
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
