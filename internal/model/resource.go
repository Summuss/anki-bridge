package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
