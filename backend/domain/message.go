package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ColelctionMessage string = "message"

type Message struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Sender     primitive.ObjectID `json:"sender" bson:"sender"`
	Content    string             `json:"content" bson:"content"`
	Chat       primitive.ObjectID `json:"chat" bson:"chat"`
	IsEdited   bool               `json:"isedited" bson:"isedited"`
	Created_at time.Time          `json:"created_at" bson:"created_at"`
	Updated_at time.Time          `json:"updated_at" bson:"updated_at"`
}

type MessageRepository interface {
	Create(context.Context, Message) (primitive.ObjectID, error)
	FetchById(context.Context, primitive.ObjectID) ([]bson.M, error)
	UpdateById(context.Context, primitive.D, primitive.D) error
	DeleteById(context.Context, primitive.ObjectID) error
}

type SendMessageUseCase interface {
	Create(context.Context, Message) (primitive.ObjectID, error)
	UpdateByFilter(context.Context, primitive.D, primitive.D) error
	FetchById(context.Context, primitive.ObjectID) ([]bson.M, error)
}

type GetMessagesUseCase interface {
	FetchById(context.Context, primitive.ObjectID) ([]bson.M, error)
}
