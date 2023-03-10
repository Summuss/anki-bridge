package model

import (
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JPSentence struct {
	BaseModel   `json:",inline" bson:",inline"`
	Sentence    string     `json:"sentence" bson:"sentence"`
	Explanation string     `json:"explanation" bson:"explanation"`
	Addition    string     `json:"addition" bson:"addition"`
	JPWords     *[]*JPWord `json:"jp_words" bson:"jp_words"`
}

func (j *JPSentence) CollectionName() string {
	return "jp_sentences"
}
func (j *JPSentence) Desc() string {
	return j.CollectionName() + ":" + j.Sentence
}

func (j *JPSentence) Save(client *mongo.Client, dbName string) error {
	dao := GetDao(client, dbName, j)
	return dao.Save(j)

}
func (j *JPSentence) duplicationCheckQuery() interface{} {
	return bson.D{{"sentence", j.Sentence}}
}

func (m *JPSentence) GetResources() *[]Resource {
	if m.GetNoteInfo().Name == common.NoteType_JPSentences_Voice_Name {
		return m.resources
	}
	resources := lo.FlatMap(
		*m.JPWords, func(item *JPWord, index int) []Resource {
			return *item.resources
		},
	)
	return &resources
}
