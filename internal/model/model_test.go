package model

type UserModel struct {
	Model `json:",inline" bson:",inline"`
	Name  string `json:"name" bson:"name" required:"true" minLen:"2"`
	Age   int    `json:"age" bson:"age"`
}

func (j *UserModel) collectionName() string {
	return "user"
}

type AddrModel struct {
	Model `json:",inline" bson:",inline"`
	State string `json:"state" bson:"state"`
	City  string `json:"city" bson:"city"`
}

func (j *AddrModel) collectionName() string {
	return "user"
}
