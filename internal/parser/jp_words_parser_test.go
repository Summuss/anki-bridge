package parser

import (
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
