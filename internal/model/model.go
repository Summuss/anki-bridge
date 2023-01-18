package model

import (
	"github.com/summuss/anki-bridge/internal/common"
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

	SetNoteTypeName(noteTypeName common.NoteTypeName)
	GetNoteTypeName() common.NoteTypeName

	// set after middle parse finished
	SetMiddleInfo(info interface{})
	GetMiddleInfo() interface{}
	// to avoid cycle import, use interface{} instead of parser.iParser
	GetParser() interface{}
	SetParser(interface{})
	GetNoteInfo() *common.NoteInfo
	SetNoteInfo(*common.NoteInfo)

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
	NoteType    common.NoteTypeName   `json:"note_type" bson:"note_type"`
	Resources   *[]primitive.ObjectID `json:"resources" bson:"resources"`
	resources   *[]Resource
	noteInfo    *common.NoteInfo
	parser      interface{}
	middleInfo  interface{}
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
func (m *BaseModel) SetNoteTypeName(noteTypeName common.NoteTypeName) {
	m.NoteType = noteTypeName
}
func (m *BaseModel) GetNoteTypeName() common.NoteTypeName {
	return m.NoteType
}

func (m *BaseModel) SetMiddleInfo(info interface{}) {
	m.middleInfo = info
}
func (m *BaseModel) GetMiddleInfo() interface{} {
	return m.middleInfo
}

func (m *BaseModel) SetNoteInfo(n *common.NoteInfo) {
	m.noteInfo = n
}

func (m *BaseModel) GetNoteInfo() *common.NoteInfo {
	return m.noteInfo
}
func (m *BaseModel) GetParser() interface{} {
	return m.parser
}

func (m *BaseModel) SetParser(p interface{}) {
	m.parser = p
}

func getZeroModel() BaseModel {
	return BaseModel{ID: primitive.NilObjectID, CreatedTime: 0, UpdateTime: 0}
}
