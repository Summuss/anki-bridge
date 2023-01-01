package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sync"
	"time"
)

var DB_NAME = "test"
var _mongoClient *mongo.Client
var _once = sync.Once{}

type IDao[T IModel] interface {
	FindById(primitive.ObjectID) (*T, error)
	FindMany(query interface{}) (*[]T, error)
	Save(model T, models ...T) error
	Delete(model T, models ...T) error
}

type Dao[T IModel] struct {
	Client *mongo.Client
}

func (d *Dao[T]) FindById(id primitive.ObjectID) (T, error) {
	var t T // T should be a pointer type which point to a concrete model struct
	tt := reflect.TypeOf(t)
	one := d.getCollection().FindOne(
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
	cursor, err := d.getCollection().Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	var res []T
	err = cursor.All(context.TODO(), &res)
	return &res, err
}
func (d *Dao[T]) Save(model T, models ...T) error {
	var insertMs []interface{}
	var updateMs []interface{}
	addToGroup := func(m T) {
		if m.GetID().IsZero() {
			insertMs = append(insertMs, m)
		} else {
			updateMs = append(updateMs, m)
		}
	}
	addToGroup(model)
	for _, m := range models {
		addToGroup(m)
	}

	ctx := context.TODO()
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		now := primitive.NewDateTimeFromTime(time.Now())
		if len(insertMs) > 0 {
			for _, m := range insertMs {
				m.(T).SetCreatedTime(now)
				m.(T).SetUpdateTime(now)
			}
			res, err := d.getCollection().InsertMany(ctx, insertMs)
			if err != nil {
				return nil, fmt.Errorf("insertMany failed, error:\n%s", err.Error())
			}
			for i, id := range res.InsertedIDs {
				objectID := id.(primitive.ObjectID)
				insertMs[i].(T).SetID(objectID)
			}

		}
		for _, m := range updateMs {
			id := m.(T).GetID()
			m.(T).SetUpdateTime(now)
			_, err := d.getCollection().ReplaceOne(ctx, bson.D{{"_id", id}}, m)
			if err != nil {
				return nil, fmt.Errorf(
					"update id=%s failed,error:\n%s", id.String(), err.Error(),
				)
			}
		}
		return nil, nil
	}
	session, err := d.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	res, err := session.WithTransaction(ctx, callback)
	println(res)
	return nil
}

func (d *Dao[T]) Delete(model T, models ...T) error {
	return nil
}
func (d *Dao[T]) getCollection() *mongo.Collection {
	var t T
	collectionName := t.collectionName()
	if len(collectionName) == 0 {
		panic(fmt.Sprintf("collectionName is empty for model %s", reflect.TypeOf(t)))
	}
	return d.Client.Database(DB_NAME).Collection(collectionName)
}

func getClient() *mongo.Client {
	_once.Do(
		func() {
			if _mongoClient != nil {
				return
			}
			uri := fmt.Sprintf(
				"mongodb://mongoadmin:secret@daemon:27017/%s?authSource=admin", DB_NAME,
			)
			var err error
			_mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
			if err != nil {
				panic(err)
			}
		},
	)
	return _mongoClient
}

func GetDao[T IModel](t T) Dao[T] {
	return Dao[T]{Client: getClient()}
}

func Disconnect() error {
	return getClient().Disconnect(context.TODO())
}

func ObjectIDFromHex(idStr string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
