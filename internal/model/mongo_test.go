package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func Test_connectToMongo(t *testing.T) {
	/*	uri := "mongodb://mongoadmin:secret@daemon:27017/test?authSource=admin"
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}
	*/
	var d = GetDao(&UserModel{})
	objectID, err := primitive.ObjectIDFromHex("63aee359aa00acdf6ffc8769")
	if err != nil {
		panic(err)
	}
	res, err := d.FindById(objectID)
	if err != nil {
		panic(err)
	}
	print(res)
}

func TestDao_FindById(t *testing.T) {
	GetDao(&UserModel{})
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
			d:       GetDao(&UserModel{}),
			args:    args{ObjectIDFromHex("63aee359aa00acdf6ffc8769")},
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
		wantErr bool
	}
	tests := []testCase[*UserModel]{
		{name: "1", d: GetDao(&UserModel{}), args: args{bson.D{}}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.d.FindMany(tt.args.query)
				if (err != nil) != tt.wantErr {
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
