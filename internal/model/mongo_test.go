package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

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

func TestDao_Save(t *testing.T) {
	model1 := UserModel{
		getZeroModel(), "liu", 27,
		&[]AddrModel{{getZeroModel(), "jp", "tokyo"}},
	}
	model2 := model1
	model2.Age = 24
	dao := GetDao(&UserModel{})
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
