package anki

import (
	"io"
	"net/http"
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

func Test_Downlaod_data(t *testing.T) {
	resp, err := http.Get("https://cache-a.oddcast.com/c_fs/8dfc5794d25cab8a28b2874f8b461044.mp3?engine=3&language=12&voice=4&text=%22%E5%89%B0%E3%81%88%22&useUTF8=1")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	println(body)

}

func TestGetAllDecks(t *testing.T) {
	decks, err := GetAllDecks()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(decks) == 0 {
		t.Errorf("desk not found")
	}

}

func TestGetAllAnkiModels(t *testing.T) {
	models, err := GetAllAnkiModels()
	if err != nil {
		t.Error(err)
		return

	}
	if len(models) == 0 {
		t.Error("size=0")
	}
}
