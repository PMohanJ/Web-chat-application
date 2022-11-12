package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Sender  primitive.ObjectID   `bson:"_id"`
	Content string               `bson:"content"`
	Chat    primitive.ObjectID   `bson:"chat"`
	ReadBy  []primitive.ObjectID `bson:"readBy"`
}
