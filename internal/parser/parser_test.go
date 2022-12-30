package parser

import (
	"reflect"
	"testing"
)

func Test_splitByNoteType(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				content: `
Word
	W1
		1-1
		1-2
	W2
		2-1

Sentence
	S1
		1-1
		1-2
	S2
		2-1
Word
	W3
		3-1
		3-2
	W4
		4-1
`,
			},
			want: map[string]string{
				"Word": `
W1
	1-1
	1-2
W2
	2-1
W3
	3-1
	3-2
W4
	4-1`,

				"Sentence": `
S1
	1-1
	1-2
S2
	2-1`,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				content: `
Word
Sentence
Word
	1-1
`,
			},
			want: map[string]string{
				"Word": `
1-1`,
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				content: `
 1234
`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "note name contain space",
			args: args{
				content: `
Word
	W1
		1-1
		1-2
	W2
		2-1

Sentence xxx
	S1
		1-1
		1-2
	S2
		2-1
Word
	W3
		3-1
		3-2
	W4
		4-1
`,
			},
			want: map[string]string{
				"Word": `
W1
	1-1
	1-2
W2
	2-1
W3
	3-1
	3-2
W4
	4-1`,

				"Sentence xxx": `
S1
	1-1
	1-2
S2
	2-1`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := splitByNoteType(tt.args.content)
				if (err != nil) != tt.wantErr {
					t.Errorf("splitByNoteType() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("splitByNoteType() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
