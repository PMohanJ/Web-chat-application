package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Email    string             `baon:"email"`
	Password string             `bson:"password"`
	Pic      string             `bson:"pic"`
	IsAdmin  bool               `bson:"isAdmin"`
}
