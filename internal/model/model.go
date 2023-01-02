package model

import (
	"go.mongodb.org/mongo-driver/bson"
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

func getZeroModel() BaseModel {
	return BaseModel{ID: primitive.NilObjectID, CreatedTime: 0, UpdateTime: 0}
}

type Resource struct {
	Metadata ResourceMetadata   `json:"metadata" bson:"metadata"`
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Length   int                `json:"length" bson:"length"`
	data     []byte
}

func (r *Resource) SetData(data []byte) {
	r.data = data
}
func (r *Resource) GetData() []byte {
	return r.data
}

type ResourceMetadata struct {
	ExtName      string             `json:"ext_name" bson:"ext_name"`
	ResourceType ResourceType       `json:"file_type" bson:"file_type"`
	Collection   string             `json:"collection" bson:"collection"`
	OwnerID      primitive.ObjectID `json:"model_id" bson:"model_id"`
	FileName     string             `json:"file_name" bson:"file_name"`
}

func (r *Resource) toBsonM() bson.M {
	return bson.M{
		"ext_name": r.Metadata.ExtName, "file_type": r.Metadata.ResourceType,
		"file_name": r.Metadata.FileName, "collection": r.Metadata.Collection,
		"model_id": r.Metadata.OwnerID,
	}
}

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

type JPSentence struct {
	BaseModel   `json:",inline" bson:",inline"`
	Sentence    string     `json:"sentence" bson:"sentence"`
	Explanation string     `json:"explanation" bson:"explanation"`
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
