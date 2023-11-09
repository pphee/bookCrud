package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Student struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	Age       int                `bson:"age"`
}
