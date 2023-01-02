package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResourceType string

var (
	Sound ResourceType = "sound"
	Video ResourceType = "video"
	Image ResourceType = "image"
)

type IModel interface {

	// for BaseModel Impl
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

	SetAnkiNoteId(int64)

	SetNoteType(noteType string)
	GetNoteType() string

	// for Concrete Model Impl
	CollectionName() string
	Save(client *mongo.Client, dbName string) error
	Desc() string
	duplicationCheckQuery() interface{}
}

type BaseModel struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	CreatedTime primitive.DateTime    `json:"created_time" bson:"created_time"`
	UpdateTime  primitive.DateTime    `json:"update_time" bson:"update_time"`
	ChangeFlag  string                `json:"change_flag"  bson:"change_flag"`
	AnkiNoteId  int64                 `json:"anki_note_id"  bson:"anki_note_id"`
	NoteType    string                `json:"note_type"bson:"note_type"`
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

func (m *BaseModel) SetAnkiNoteId(id int64) {
	m.AnkiNoteId = id
}
func (m *BaseModel) SetNoteType(noteType string) {
	m.NoteType = noteType
}
func (m *BaseModel) GetNoteType() string {
	return m.NoteType
}

func getZeroModel() BaseModel {
	return BaseModel{ID: primitive.NilObjectID, CreatedTime: 0, UpdateTime: 0}
}
