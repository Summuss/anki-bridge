package common

import "testing"

func TestAddExtra(t *testing.T) {
	noateInfo := &NoteInfo{}
	SetExtra(noateInfo, NO_JPWORD_TTS, "world")
	s := ""
	FetchExtraByKey(noateInfo, NO_JPWORD_TTS, &s)
	if s != "world" {
		t.Errorf("unexpected value")
	}
}
