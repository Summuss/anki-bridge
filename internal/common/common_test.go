package common

import (
	"testing"
)

func Test_curlGetData(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *[]byte
		wantErr bool
	}{
		{
			name: "1",
			args: args{url: `https://cache-a.oddcast.com/c_fs/c3fc3fe119f11f4e33d9a88322cf363a.mp3?engine=3&language=12&voice=4&text=%E6%97%A5%E6%9C%AC%E8%AA%9E&useUTF8=1`},
			want: nil, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := CurlGetData(tt.args.url)
				if (err != nil) != tt.wantErr {
					t.Errorf("CurlGetData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got == nil {
					t.Errorf("CurlGetData() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
