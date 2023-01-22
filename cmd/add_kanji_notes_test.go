package cmd

import "testing"

func Test_addKanji(t *testing.T) {
	err := addKanji()
	if err != nil {
		t.Errorf(err.Error())
	}
}
