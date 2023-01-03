package common

type NoteType string

var (
	NoteType_JPWords       NoteType = "Jp Words"
	NoteType_JPSentences   NoteType = "Jp Sentences"
	NoteType_JPRecognition NoteType = "认识"
	NoteTypeList                    = []NoteType{
		NoteType_JPWords, NoteType_JPRecognition, NoteType_JPSentences,
	}
)