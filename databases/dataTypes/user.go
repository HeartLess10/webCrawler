package dataTypes

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserDataType struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
	Age  int                `bson:"age,omitempty"`
	Job  string             `bson:"job,omitempty"`
}
