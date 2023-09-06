package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var JpListenCollectionName = "jp_listen"

type JpListen[T any] struct {
	BaseModel   `json:",inline" bson:",inline"`
	JpListenKey string `json:"jp_listen_key" bson:"jp_listen_key"`
	ExtInfo     *T     `json:"ext_info" bson:"ext_info"`
}

func (j *JpListen[T]) CollectionName() string {
	return JpListenCollectionName
}

func (j *JpListen[T]) Save(client *mongo.Client, dbName string) error {
	dao := GetDao(client, dbName, j)
	return dao.Save(j)
}

func (j *JpListen[T]) Desc() string {
	return j.CollectionName() + ":" + j.JpListenKey
}

func (j *JpListen[T]) duplicationCheckQuery() interface{} {
	return bson.D{{"jp_listen_key", j.JpListenKey}}
}
