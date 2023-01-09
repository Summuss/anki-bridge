package common

type NoteInfo struct {
	Name          NoteTypeName
	Title         string `yaml:"title"`
	Desk          string `yaml:"desk"`
	AnkiNoteModel string `yaml:"anki-note-model"`
}

type NoteTypeName string

var (
	NoteType_JPWords_Name       NoteTypeName = "JPWords"
	NoteType_JPSentences_Name   NoteTypeName = "JPSentences"
	NoteType_JPRecognition_Name NoteTypeName = "JPRecognition"
	NoteType_Kanji_Name         NoteTypeName = "Kanji"

	NoteTypeNameList []NoteTypeName = []NoteTypeName{
		NoteType_JPWords_Name,
		NoteType_JPSentences_Name,
		NoteType_JPRecognition_Name,
		NoteType_Kanji_Name,
	}
)
