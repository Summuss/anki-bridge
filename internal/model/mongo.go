package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type IDao[T IModel] interface {
	FindById(primitive.ObjectID) (*T, error)
	FindMany(interface{}) (*[]T, error)
}

type Dao[T IModel] struct {
	Client *mongo.Client
}

func (d *Dao[T]) FindById(id primitive.ObjectID) (T, error) {
	var t T // T should be a pointer type which point to a concrete model struct
	tt := reflect.TypeOf(t)

	collectionName := t.collectionName()
	if len(collectionName) == 0 {
		panic(fmt.Sprintf("collectionName is empty for model %s", tt.Name()))
	}
	one := d.Client.Database("test").Collection(collectionName).FindOne(
		context.TODO(),
		bson.D{{"_id", id}},
	)
	if one.Err() != nil {
		return t, one.Err()
	}
	t = reflect.New(tt.Elem()).Interface().(T)
	err := one.Decode(t)
	return t, err
}

func (d *Dao[T]) FindMany(query interface{}) (*[]T, error) {
	var t T
	tt := reflect.TypeOf(t)

	collectionName := t.collectionName()
	if len(collectionName) == 0 {
		panic(fmt.Sprintf("collectionName is empty for model %s", tt.Name()))
	}
	cursor, err := d.Client.
		Database("test").
		Collection(collectionName).
		Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	var res []T
	err = cursor.All(context.TODO(), &res)
	return &res, err
}

func GetDao[T IModel](t T) Dao[T] {
	uri := "mongodb://mongoadmin:secret@daemon:27017/test?authSource=admin"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return Dao[T]{Client: client}

}

func ObjectIDFromHex(idStr string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
