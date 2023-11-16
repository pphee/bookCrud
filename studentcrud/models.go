package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Student struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstName" faker:"first_name"`
	LastName  string             `bson:"lastName" faker:"last_name"`
	Age       int                `bson:"age" faker:"boundary_start=18, boundary_end=30"`
}
