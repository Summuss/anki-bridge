package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

var TestDB = "test"
var mongoClient *mongo.Client

func init() {
	uri := "mongodb://mongoadmin:secret@daemon:27017/test?authSource=admin"
	var err error
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
}

func TestDao_FindById(t *testing.T) {
	GetDao(mongoClient, TestDB, &UserModel{})
	type args struct {
		id primitive.ObjectID
	}
	type testCase[T IModel] struct {
		name    string
		d       Dao[T]
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[*UserModel]{
		{
			name:    "1",
			d:       GetDao(mongoClient, TestDB, &UserModel{}),
			args:    args{ObjectIDFromHex("63b1155b63ac6ba5560e0f80")},
			want:    &UserModel{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.d.FindById(tt.args.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("FindById() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got.GetID() != tt.args.id {
					t.Errorf("not found")
				}
			},
		)
	}

	jpWordDao := GetDao(mongoClient, TestDB, &JPWord{})
	res, err := jpWordDao.FindById(ObjectIDFromHex("6180a5d05c1e8d3bb3362f3f"))
	if err != nil {
		t.Errorf(err.Error())
	}
	println(res.CreatedTime.Time().String())
}

func TestDao_FindMany(t *testing.T) {
	type args struct {
		query interface{}
	}
	type testCase[T IModel] struct {
		name    string
		d       Dao[T]
		args    args
		want    *[]T
		wantErr int
	}
	tests := []testCase[*JPWord]{
		{
			name: "1", d: GetDao(mongoClient, TestDB, &JPWord{}), args: args{bson.D{}},
			want:    nil,
			wantErr: 1,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.d.FindMany(tt.args.query)
				if (err != nil) != false {
					t.Errorf("FindMany() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got == nil || len(*got) == 0 {
					t.Errorf("not found")
					return
				}
			},
		)
	}
}

func TestDao_Save(t *testing.T) {
	model1 := UserModel{
		getZeroModel(), "liu", 27,
		&[]AddrModel{{getZeroModel(), "jp", "tokyo"}},
	}
	model2 := model1
	model2.Age = 24
	dao := GetDao(mongoClient, TestDB, &UserModel{})
	err := dao.Save(&model1)
	if err != nil {
		t.Errorf(err.Error())
	}
	model1.Age = 26
	err = dao.Save(&model1, &model2)
	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestDao_loadResources(t *testing.T) {
	type args[T IModel] struct {
		t T
	}
	type testCase[T IModel] struct {
		name    string
		d       Dao[T]
		args    args[T]
		wantErr bool
	}

	dao := GetDao(mongoClient, TestDB, &JPWord{})
	jpWord := JPWord{}
	jpWord.Resources = &[]primitive.ObjectID{
		ObjectIDFromHex("6180a5b55c1e8d3bb3362f36"),
		ObjectIDFromHex("6180a5ba5c1e8d3bb3362f3d"),
		ObjectIDFromHex("6180a5d05c1e8d3bb3362f40"),
	}
	tests := []testCase[*JPWord]{
		{
			name:    "1",
			d:       dao,
			args:    args[*JPWord]{&jpWord},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := tt.d.loadResources(tt.args.t); (err != nil) != tt.wantErr {
					t.Errorf("loadResources() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
		if len(*tt.args.t.resources) == 0 {
			t.Errorf("empty resources")
		}
		resource := (*tt.args.t.resources)[0]
		if len(resource.data) == 0 {
			t.Errorf("empty data")

		} else {
			fo, err := os.Create(resource.FileName)
			if err != nil {
				panic(err)
			}
			_, err = fo.Write(resource.data)
			if err != nil {
				panic(err)
			}

		}
	}
}
