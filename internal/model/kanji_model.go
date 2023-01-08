package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var KanjiCollectionName = "kanji"

type Kanji struct {
	BaseModel `json:",inline" bson:",inline"`
	Kanji     string        `json:"kanji"  bson:"kanji"`
	Prons     *[]*KanjiPron `json:"prons" bson:"prons"`
}

type KanjiPron struct {
	Pron    string `json:"pron" bson:"pron"`
	Example string `json:"example" bson:"example"`
}

func (j *Kanji) CollectionName() string {
	return KanjiCollectionName
}

func (j *Kanji) Save(client *mongo.Client, dbName string) error {
	dao := GetDao(client, dbName, j)
	return dao.Save(j)
}

func (j *Kanji) Desc() string {
	return j.CollectionName() + ":" + j.Kanji
}

func (j *Kanji) duplicationCheckQuery() interface{} {
	return bson.D{{"kanji", j.Kanji}}
}
