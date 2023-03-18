package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JPCommonNote struct {
	BaseModel `json:",inline" bson:",inline"`
	Front     string `json:"front" bson:"front"`
	Back      string `json:"back" bson:"back"`
}

func (j *JPCommonNote) CollectionName() string {
	return "jp_common_notes"
}
func (j *JPCommonNote) Desc() string {
	return j.CollectionName() + ":" + j.Front
}

func (j *JPCommonNote) Save(client *mongo.Client, dbName string) error {
	dao := GetDao(client, dbName, j)
	return dao.Save(j)

}
func (j *JPCommonNote) duplicationCheckQuery() interface{} {
	return bson.D{{"front", j.Front}, {"back", j.Back}}
}

func (m *JPCommonNote) GetResources() *[]Resource {
	return &[]Resource{}
}
