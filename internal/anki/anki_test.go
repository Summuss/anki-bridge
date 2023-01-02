package anki

import (
	"reflect"
	"testing"
)

func Test_requestAnki(t *testing.T) {
	type args struct {
		action string
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				action: "notesInfo",
				params: map[string]interface{}{"notes": []int64{1671279148859}},
			}, want: nil, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := requestAnki(tt.args.action, tt.args.params)
				if (err != nil) != tt.wantErr {
					t.Errorf("requestAnki() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("requestAnki() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
