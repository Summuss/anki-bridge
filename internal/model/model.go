package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IModel interface {
	collectionName() string
	SetID(id primitive.ObjectID)
	GetID() primitive.ObjectID
	GetCreatedTime() primitive.DateTime
	GetUpdateTime() primitive.DateTime
	SetCreatedTime(primitive.DateTime)
	SetUpdateTime(primitive.DateTime)
}

type Model struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedTime primitive.DateTime `json:"created_time" bson:"created_time"`
	UpdateTime  primitive.DateTime `json:"update_time" bson:"update_time"`
}

func (m *Model) GetID() primitive.ObjectID {
	return m.ID
}
func (m *Model) SetID(id primitive.ObjectID) {
	m.ID = id
}

func (m *Model) GetCreatedTime() primitive.DateTime {
	return m.CreatedTime
}

func (m *Model) GetUpdateTime() primitive.DateTime {
	return m.GetUpdateTime()
}
func (m *Model) SetUpdateTime(time primitive.DateTime) {
	m.UpdateTime = time
}
func (m *Model) SetCreatedTime(time primitive.DateTime) {
	m.CreatedTime = time
}

func getZeroModel() Model {
	return Model{ID: primitive.NilObjectID, CreatedTime: 0, UpdateTime: 0}
}

type FileModel struct {
}

type JPWord struct {
	Model       `json:",inline" bson:",inline"`
	AnkiNoteId  string               `json:"anki_note_id"  bson:"anki_note_id"`
	ChangeFlag  string               `json:"change_flag"  bson:"change_flag"`
	Hiragana    string               `json:"hiragana"  bson:"hiragana"`
	Mean        string               `json:"mean"  bson:"mean"`
	Pitch       string               `json:"pitch"  bson:"pitch"`
	Resources   []primitive.ObjectID `json:"resources"  bson:"resources"`
	Spell       string               `json:"spell"  bson:"spell"`
	WordClasses []int                `json:"word_classes"  bson:"word_classes"`
}

func (j *JPWord) collectionName() string {
	return "jp_words"
}
