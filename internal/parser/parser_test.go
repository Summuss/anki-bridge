package parser

import (
	"github.com/summuss/anki-bridge/internal/model"
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
- [[Word]]
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

func TestCheckInput(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				txt: `
- [[Jp Sentences]]
	- ‎原因不明の病に {{cloze 冒される}} ので
		- 冒す
			- 侵袭；患（病）。（害を及ぼす。）
			- おかす　2　7
	- 前職で重いヘルニアを {{cloze ‎患}} ってらしたみたいで
		- 患う
			- 患病，生病。〔病気になる。〕
			- わずらう　0　7
- [[Jp Words]]
	- 出来事
		- (偶发)的事件，变故。（持ち上がった事件・事柄。）
		- できごと　2　１
	- せっかち
		- 性急；急躁（落ち着きがなく，先へ先へと急ぐ・ことさま。また，そのような性質の人。性急）。
		- せっかち　1　3

`,
			},
			wantErr: false,
		},
		{
			name: "structure error",
			args: args{
				txt: `
- [[Jp Sentences]]
	- ‎原因不明の病に {{cloze 冒される}} ので
		- 冒す
			- 侵袭；患（病）。（害を及ぼす。）
			- おかす　2　7
	- 前職で重いヘルニアを {{cloze ‎患}} ってらしたみたいで
		- 患う
			- 患病，生病。〔病気になる。〕
			- わずらう　0　7
- [[Jp Words]]
 abc
`,
			},
			wantErr: true,
		},

		{
			name: "note error1",
			args: args{
				txt: `
- [[Jp Sentences]]
	- ‎原因不明の病に {{cloze 冒される}} ので
		- 冒す
			- 侵袭；患（病）。（害を及ぼす。）
			- おかす　2　7
	- 前職で重いヘルニアを {{cloze ‎患}} ってらしたみたいで
		- 患う
			- 患病，生病。〔病気になる。〕
			- わずらう　0　7
- [[Jp Words]]
	- 出来事
		- (偶发)的事件，变故。（持ち上がった事件・事柄。）
		- できごと　2　X
	- せっかち
		- 性急；急躁（落ち着きがなく，先へ先へと急ぐ・ことさま。また，そのような性質の人。性急）。
		- せっかち　1　3
`,
			},
			wantErr: true,
		},
		{
			name: "note error2",
			args: args{
				txt: `
- [[Jp Sentences]]
	- ‎原因不明の病に {{cloze 冒される}} ので
		- 冒す
			- 侵袭；患（病）。（害を及ぼす。）
			- おかす　2　0
	- 前職で重いヘルニアを {{cloze ‎患}} ってらしたみたいで
		- 患う
			- 患病，生病。〔病気になる。〕
			- わずらう　0　7
- [[Jp Words]]
	- 出来事
		- (偶发)的事件，变故。（持ち上がった事件・事柄。）
		- できごと　2　X
	- せっかち
		- 性急；急躁（落ち着きがなく，先へ先へと急ぐ・ことさま。また，そのような性質の人。性急）。
		- せっかち　1　3
`,
			},
			wantErr: true,
		},
		{
			name: "empty note",
			args: args{
				txt: `
- [[Jp Sentences]]
- [[Jp Words]]
	- 出来事
		- (偶发)的事件，变故。（持ち上がった事件・事柄。）
		- できごと　2　A
	- せっかち
		- 性急；急躁（落ち着きがなく，先へ先へと急ぐ・ことさま。また，そのような性質の人。性急）。
		- せっかち　1　3
`,
			},
			wantErr: false,
		},
		{
			name: "unknown note",
			args: args{
				txt: `
- [[UNKNOWN]]
	- xxx
- [[Jp Words]]
	- 出来事
		- (偶发)的事件，变故。（持ち上がった事件・事柄。）
		- できごと　2　A
	- せっかち
		- 性急；急躁（落ち着きがなく，先へ先へと急ぐ・ことさま。また，そのような性質の人。性急）。
		- せっかち　1　3
`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := CheckInput(tt.args.txt)
				if (err != nil) != tt.wantErr {
					t.Errorf("CheckInput() error = %v, wantErr %v", err, tt.wantErr)
				}
				if err != nil {
					println(err.Error())
				}
			},
		)
	}
}

func Test_computeNoteName(t *testing.T) {
	type args struct {
		rowNoteName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1", args: args{rowNoteName: "- [[Jp Sentences]]"}, want: "Jp Sentences"},
		{name: "2", args: args{rowNoteName: "- [[Jp Sentences]]  "}, want: "Jp Sentences"},
		{name: "3", args: args{rowNoteName: "JXX"}, want: "JXX"},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := computeRealNoteName(tt.args.rowNoteName); got != tt.want {
					t.Errorf("computeRealNoteName() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestParse(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    *[]model.IModel
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				text: `
- [[Jp Sentences]]
	- ‎原因不明の病に {{cloze 冒される}} ので
		- 冒す
			- 侵袭；患（病）。（害を及ぼす。）
			- おかす　2　7
	- 前職で重いヘルニアを {{cloze ‎患}} ってらしたみたいで
		- 患う
			- 患病，生病。〔病気になる。〕
			- わずらう　0　7
- [[Jp Words]]
	- 出来事
		- (偶发)的事件，变故。（持ち上がった事件・事柄。）
		- できごと　2　１
	- せっかち
		- 性急；急躁（落ち着きがなく，先へ先へと急ぐ・ことさま。また，そのような性質の人。性急）。
		- せっかち　1　3

`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := Parse(tt.args.text)
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
