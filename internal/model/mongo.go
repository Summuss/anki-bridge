package model

import (
	"context"
	"fmt"
	"github.com/summuss/anki-bridge/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"reflect"
	"time"
)

type IDao[T IModel] interface {
	FindById(primitive.ObjectID) (*T, error)
	FindMany(query interface{}) (*[]T, error)
	Save(model T, models ...T) error
	Delete(model T, models ...T) error
}

type Dao[T IModel] struct {
	Client *mongo.Client
	DBName string
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
	if err != nil {
		return t, err
	}
	err = d.loadResources(t)
	return t, err
}

func (d *Dao[T]) FindMany(query interface{}) (*[]T, error) {
	cursor, err := d.getCollection().Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	var res []T
	err = cursor.All(context.TODO(), &res)
	if err != nil {
		return nil, err
	}
	err = util.DoParallel(
		&res, func(t *T) error {
			err = d.loadResources(*t)
			if err != nil {
				return fmt.Errorf("load resources failed,error:\n%s", err.Error())
			}
			return nil
		},
	)
	return &res, err
}
func (d *Dao[T]) loadResources(t T) error {
	if t.GetResourceIDs() == nil {
		return nil
	}
	rsiSize := len(*t.GetResourceIDs())
	if rsiSize > 0 {
		db := d.Client.Database(d.DBName)
		bucket, err := gridfs.NewBucket(db)
		if err != nil {
			return err
		}
		cursor, err := bucket.Find(bson.D{{"_id", bson.D{{"$in", t.GetResourceIDs()}}}})
		if err != nil {
			return err
		}
		var rs []Resource
		err = cursor.All(context.TODO(), &rs)
		if err != nil {
			return err
		}
		if len(rs) != rsiSize {
			fmt.Printf("warnning: resource id's num:%d ,only %d found", rsiSize, len(rs))
		}
		err = util.DoParallel(
			&rs,
			func(r *Resource) error {
				if r.Length == 0 {
					return nil
				}
				downloadStream, err := bucket.OpenDownloadStream(r.Id)
				if err != nil {
					return fmt.Errorf(
						"download file %s(id=%s) failed,error:\n%s", r.FileName, r.Id.String(),
						err.Error(),
					)
				}
				r.data = make([]byte, r.Length)
				_, err = downloadStream.Read(r.data)
				if err != nil {
					return fmt.Errorf(
						"download file %s(id=%s) failed,error:\n%s", r.FileName, r.Id.String(),
						err.Error(),
					)
				}
				return nil
			},
		)
		if err != nil {
			return err
		}
		t.SetResources(&rs)
	}
	return nil
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
	_, err = session.WithTransaction(ctx, callback)
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
	return d.Client.Database(d.DBName).Collection(collectionName)
}

func GetDao[T IModel](client *mongo.Client, dbName string, t T) Dao[T] {
	return Dao[T]{Client: client, DBName: dbName}
}

func ObjectIDFromHex(idStr string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
