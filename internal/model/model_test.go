package model

import "go.mongodb.org/mongo-driver/mongo"

type UserModel struct {
	BaseModel `json:",inline" bson:",inline"`
	Name      string       `json:"name" bson:"name" required:"true" minLen:"2"`
	Age       int          `json:"age" bson:"age"`
	Addr      *[]AddrModel `json:"addr" bson:"addr"`
}

func (j *UserModel) CollectionName() string {
	return "user"
}
func (j *UserModel) Save(client *mongo.Client, dbName string) error {
	dao := GetDao(client, dbName, j)
	return dao.Save(j)
}

func (j *UserModel) Desc() string {
	return j.CollectionName() + ":" + j.Name
}

type AddrModel struct {
	b     BaseModel `json:",inline" bson:",inline"`
	State string    `json:"state" bson:"state"`
	City  string    `json:"city" bson:"city"`
}

func (j *AddrModel) collectionName() string {
	return "user"
}
