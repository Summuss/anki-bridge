package splitter

import (
	"github.com/summuss/anki-bridge/internal/common"
)

var splitters *[]splitter = &[]splitter{}

type splitter interface {
	TargetNoteType() common.NoteInfo
	Split(string) (*[]string, error)
}

type simpleSplitter struct {
}

func (s simpleSplitter) Split(rawNotes string) (*[]string, error) {
	return common.SplitByNoIndentLine(rawNotes)

}

func Split(subText string, noteType common.NoteInfo) (*[]string, error) {
	for i := 0; i < len(*splitters); i++ {
		if (*splitters)[i].TargetNoteType() == noteType {
			return (*splitters)[i].Split(subText)
		}
	}
	return simpleSplitter{}.Split(subText)
}
