package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JPWord struct {
	BaseModel   `json:",inline" bson:",inline"`
	Hiragana    string `json:"hiragana"  bson:"hiragana"`
	Mean        string `json:"mean"  bson:"mean"`
	Pitch       string `json:"pitch"  bson:"pitch"`
	Spell       string `json:"spell"  bson:"spell"`
	WordClasses []int  `json:"word_classes"  bson:"word_classes"`
}

func (j *JPWord) CollectionName() string {
	return "jp_words"
}

func (j *JPWord) Save(client *mongo.Client, dbName string) error {
	dao := GetDao(client, dbName, j)
	return dao.Save(j)
}

func (j *JPWord) Desc() string {
	return j.CollectionName() + ":" + j.Spell
}

func (j *JPWord) duplicationCheckQuery() interface{} {
	return bson.D{{"spell", j.Spell}, {"mean", j.Mean}}
}
