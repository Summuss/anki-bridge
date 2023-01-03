package model

import (
	"bytes"
	"context"
	"fmt"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
	"time"
)

var MongoClient *mongo.Client

func init() {
	var err error
	MongoClient, err = mongo.Connect(
		context.TODO(), options.Client().ApplyURI(config.Conf.MongoConnectURL),
	)
	if err != nil {
		panic(err)
	}
}

type Dao[T IModel] struct {
	Client *mongo.Client
	DBName string
}

type ExistError struct {
}

func (e ExistError) Error() string {
	return "model existed"
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
	err = common.DoParallel(
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

func (d *Dao[T]) saveResources(t T) error {
	if t.GetResources() == nil {
		return nil
	}

	if t.GetResources() == nil || len(*t.GetResources()) == 0 {
		return nil
	}
	ris := common.SafeList[primitive.ObjectID]{}
	err := common.DoParallel(
		t.GetResources(), func(r *Resource) error {
			if len(r.Metadata.FileName) == 0 {
				return fmt.Errorf("file name is empty")
			}
			if len(r.data) == 0 {
				return fmt.Errorf("data is empty, fileName:%s", r.Metadata.FileName)
			}
			r.Metadata.OwnerID = t.GetID()
			r.Length = len(r.data)
			r.Metadata.Collection = t.CollectionName()
			metadata := options.GridFSUpload().SetMetadata(r.toBsonM())
			db := d.Client.Database(d.DBName)
			bucket, err := gridfs.NewBucket(db)
			objectID, err := bucket.UploadFromStream(
				r.Metadata.FileName, bytes.NewReader(r.data),
				metadata,
			)
			if err != nil {
				return fmt.Errorf(
					"upload file %s failed,error:\n%s", r.Metadata.FileName, err.Error(),
				)
			}
			r.Id = objectID
			ris.Add(objectID)
			return nil
		},
	)
	ids := ris.ToSlice()
	t.SetResourceIDs(&ids)
	return err
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
			log.Printf("warnning: resource id's num:%d ,only %d found", rsiSize, len(rs))
		}
		err = common.DoParallel(
			&rs,
			func(r *Resource) error {
				if r.Length == 0 {
					return nil
				}
				downloadStream, err := bucket.OpenDownloadStream(r.Id)
				if err != nil {
					return fmt.Errorf(
						"download file %s(id=%s) failed,error:\n%s", r.Metadata.FileName,
						r.Id.String(),
						err.Error(),
					)
				}
				r.data = make([]byte, r.Length)
				_, err = downloadStream.Read(r.data)
				if err != nil {
					return fmt.Errorf(
						"download file %s(id=%s) failed,error:\n%s", r.Metadata.FileName,
						r.Id.String(),
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
		updateMs = append(updateMs, m)
		if m.GetID().IsZero() {
			insertMs = append(insertMs, m)
		}
	}
	addToGroup(model)
	for _, m := range models {
		addToGroup(m)
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		now := primitive.NewDateTimeFromTime(time.Now())
		if len(insertMs) > 0 {
			for _, m := range insertMs {
				m.(T).SetCreatedTime(now)
				m.(T).SetUpdateTime(now)
			}
			// transaction disabled
			for _, m := range insertMs {
				t := m.(T)
				err := d.CheckDuplication(t)
				if err != nil {
					return nil, err
				}
			}
			res, err := d.getCollection().InsertMany(context.TODO(), insertMs)
			if err != nil {
				return nil, fmt.Errorf("insertMany failed, error:\n%s", err.Error())
			}
			for i, id := range res.InsertedIDs {
				objectID := id.(primitive.ObjectID)
				insertMs[i].(T).SetID(objectID)
			}
			err = common.DoParallel(
				&insertMs, func(i *interface{}) error {
					return d.saveResources((*i).(T))
				},
			)
			if err != nil {
				return nil, err
			}

		}
		for _, m := range updateMs {
			id := m.(T).GetID()
			m.(T).SetUpdateTime(now)
			_, err := d.getCollection().ReplaceOne(context.TODO(), bson.D{{"_id", id}}, m)
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
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), callback)
	return err
}

func (d *Dao[T]) CheckDuplication(m T) error {
	count, err := d.getCollection().CountDocuments(context.TODO(), m.duplicationCheckQuery())
	if err != nil {
		return fmt.Errorf("error happend when check duplication:%s", err.Error())
	}
	if count > 0 {
		return ExistError{}
	}
	return nil
}

func (d *Dao[T]) Delete(model T, models ...T) error {
	return nil
}
func (d *Dao[T]) getCollection() *mongo.Collection {
	var t T
	collectionName := t.CollectionName()
	if len(collectionName) == 0 {
		panic(fmt.Sprintf("CollectionName is empty for model %s", reflect.TypeOf(t)))
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
