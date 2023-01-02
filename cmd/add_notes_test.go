package cmd

import "testing"

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
- [[Jp Words]]
	- 衰える
		- 而且并且
		- あまつさえ　1　１

`,
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := AddNotes(tt.args.text); (err != nil) != tt.wantErr {
					t.Errorf("AddNotes() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
