package common

type NoteInfo struct {
	Name          NoteTypeName
	Title         string `yaml:"title"`
	Desk          string `yaml:"desk"`
	AnkiNoteModel string `yaml:"anki-note-model"`
	extra         map[ExtraKey]interface{}
}

func SetExtra(noteInfo *NoteInfo, key ExtraKey, value interface{}) {
	if noteInfo.extra == nil {
		noteInfo.extra = make(map[ExtraKey]interface{})
	}
	noteInfo.extra[key] = value
}

func FetchExtraByKey[T any](noteInfo *NoteInfo, key ExtraKey, t *T) {
	if noteInfo.extra != nil {
		*t = noteInfo.extra[key].(T)
	}
}

type ExtraKey int

const (
	NO_JPWORD_TTS ExtraKey = 1 << iota
)

type NoteTypeName string

const (
	NoteType_JPWords_Name           NoteTypeName = "JPWords"
	NoteType_JPSentences_Name       NoteTypeName = "JPSentences"
	NoteType_JPRecognition_Name     NoteTypeName = "JPRecognition"
	NoteType_Kanji_Name             NoteTypeName = "Kanji"
	NoteType_JPSentences_Voice_Name NoteTypeName = "JPSentencesVoice"
	NoteType_JPCommonNotes_Name     NoteTypeName = "JPCommonNotes"
)

var NoteTypeNameList = []NoteTypeName{
	NoteType_JPWords_Name,
	NoteType_JPSentences_Name,
	NoteType_JPRecognition_Name,
	NoteType_Kanji_Name,
	NoteType_JPSentences_Voice_Name,
	NoteType_JPCommonNotes_Name,
}
