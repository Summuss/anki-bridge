package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResourceType string

var (
	Sound ResourceType = "sound"
	Video ResourceType = "video"
	Image ResourceType = "image"
)

type IModel interface {
	collectionName() string

	SetID(id primitive.ObjectID)
	GetID() primitive.ObjectID

	GetCreatedTime() primitive.DateTime
	SetCreatedTime(primitive.DateTime)

	GetUpdateTime() primitive.DateTime
	SetUpdateTime(primitive.DateTime)

	GetResources() *[]Resource
	SetResources(*[]Resource)

	GetResourceIDs() *[]primitive.ObjectID
	SetResourceIDs(*[]primitive.ObjectID)
}

type BaseModel struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	CreatedTime primitive.DateTime    `json:"created_time" bson:"created_time"`
	UpdateTime  primitive.DateTime    `json:"update_time" bson:"update_time"`
	Resources   *[]primitive.ObjectID `json:"resources" bson:"resources"`
	resources   *[]Resource
}

func (m *BaseModel) GetID() primitive.ObjectID {
	return m.ID
}
func (m *BaseModel) SetID(id primitive.ObjectID) {
	m.ID = id
}

func (m *BaseModel) GetCreatedTime() primitive.DateTime {
	return m.CreatedTime
}

func (m *BaseModel) GetUpdateTime() primitive.DateTime {
	return m.GetUpdateTime()
}
func (m *BaseModel) SetUpdateTime(time primitive.DateTime) {
	m.UpdateTime = time
}
func (m *BaseModel) SetCreatedTime(time primitive.DateTime) {
	m.CreatedTime = time
}

func (m *BaseModel) GetResources() *[]Resource {
	return m.resources
}

func (m *BaseModel) SetResources(rs *[]Resource) {
	m.resources = rs
}

func (m *BaseModel) GetResourceIDs() *[]primitive.ObjectID {
	return m.Resources
}

func (m *BaseModel) SetResourceIDs(ris *[]primitive.ObjectID) {
	m.Resources = ris
}

func getZeroModel() BaseModel {
	return BaseModel{ID: primitive.NilObjectID, CreatedTime: 0, UpdateTime: 0}
}

type Resource struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	ExtName      string             `json:"ext_name" bson:"ext_name"`
	ResourceType ResourceType       `json:"file_type" bson:"file_type"`
	FileName     string             `json:"file_name" bson:"file_name"`
	Collection   string             `json:"collection" bson:"collection"`
	OwnerID      primitive.ObjectID `json:"model_id" bson:"model_id"`
	Length       int                `json:"length" bson:"length"`
	data         []byte
}

func (r Resource) toBsonM() bson.M {
	return bson.M{
		"ext_name": r.ExtName, "file_type": r.ResourceType,
		"file_name": r.FileName, "collection": r.Collection, "model_id": r.OwnerID,
	}
}

type JPWord struct {
	BaseModel   `json:",inline" bson:",inline"`
	AnkiNoteId  int64  `json:"anki_note_id"  bson:"anki_note_id"`
	ChangeFlag  string `json:"change_flag"  bson:"change_flag"`
	Hiragana    string `json:"hiragana"  bson:"hiragana"`
	Mean        string `json:"mean"  bson:"mean"`
	Pitch       string `json:"pitch"  bson:"pitch"`
	Spell       string `json:"spell"  bson:"spell"`
	WordClasses []int  `json:"word_classes"  bson:"word_classes"`
}

func (j *JPWord) collectionName() string {
	return "jp_words"
}
