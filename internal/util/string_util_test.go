package util

import (
	"reflect"
	"testing"
)

func TestSplitByNoIndentLine(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "1", args: args{
				txt: `
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
			},
			wantErr: false,
		},
		{
			name: "2", args: args{
				txt: `
W1
	1-1
	1-2
W2 xxx
	2-1
`,
			}, want: []string{
				`W1
	1-1
	1-2
`, `W2 xxx
	2-1
`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := SplitByNoIndentLine(tt.args.txt)
				if (err != nil) != tt.wantErr {
					t.Errorf("SplitByNoIndentLine() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SplitByNoIndentLine() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
