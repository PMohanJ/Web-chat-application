package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id         primitive.ObjectID   `json:"_id,omitempty" bson:"_id"`
	Sender     primitive.ObjectID   `json:"_sender" bson:"_id"`
	Content    string               `json:"content" bson:"content"`
	Chat       primitive.ObjectID   `json:"chat" bson:"chat"`
	ReadBy     []primitive.ObjectID `json:"readBy" bson:"readBy"`
	Created_at time.Time            `json:"created_at" bson:"created_at"`
	Updated_at time.Time            `json:"updated_at" bson:"updated_at"`
}