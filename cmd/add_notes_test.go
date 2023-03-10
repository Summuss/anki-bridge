package cmd

import (
	"testing"
)

func TestAddNotes(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1", args: args{
				text: `
- [[Jp Recognition]]
	- 後回し
		- 推迟，往后推，缓办。（順番を変えてあとに遅らせること。）2
		- あとまわし　3　１
	- ニヤニヤ
		- 声を立てずに表情だけで、一人意味ありげに薄笑いを浮かべる様子。/不发出声音，一个人独笑。阴笑，不怀好意的笑。2
		- にやにや　1　4

- [[Jp Sentences]]
	- 君の一生の思い出、 {{cloze しかと}} {{cloze  見届け}} たぞ。2
		- 確と
			- 〔たしかに〕确实,明确,准确.
			- しかと　2　４
		- 見届ける
			- 看到，看准，看清；一直看到最后，结束，用眼看，确认。（最後まで見る。また、しっかり見る。）
			- みとどける　0　７
- [[Jp Words]]
	- 免れる
		- 免，避免；摆脱。（うまくさける。）2
		- まぬかれる　3　７
	- 遮る
		- 遮，遮挡，遮住，遮蔽。（光の照射や視界を邪魔する。妨げる。）2
		- さえぎる　3　７

`,
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := addNotes(tt.args.text); (err != nil) != tt.wantErr {
					t.Errorf("addNotes() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
