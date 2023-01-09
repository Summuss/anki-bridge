package config

import (
	"fmt"
	"testing"
)

func TestConfig_GetNoteInfoByTitle(t *testing.T) {
	info, _ := Conf.GetNoteInfoByTitle("Jp Words")
	fmt.Println(info)
}
