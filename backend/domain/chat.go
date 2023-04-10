package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionChat string = "chat"

type Chat struct {
	Id            primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	IsGroupChat   bool                 `json:"isGroupChat" bson:"isGroupChat"` // should default to false
	ChatName      string               `json:"chatName" bson:"chatName"`
	Users         []primitive.ObjectID `json:"users" bson:"users"`
	LatestMessage primitive.ObjectID   `json:"latestMessage" bson:"latestMessage"`
	GroupAdmin    primitive.ObjectID   `json:"groupAdmin" bson:"groupAdmin"`
	GroupPic      string               `json:"groupPic" bson:"groupPic"`
	Created_at    time.Time            `json:"created_at" bson:"created_at"`
	Updated_at    time.Time            `json:"updated_at" bson:"updated_at"`
}

func GetDefaultGroupPic() string {
	return "http://res.cloudinary.com/dkqc4za4f/image/upload/v1677908252/rf8wxtbrrxy61rvsa5hk.png"
}

type ChatRepository interface {
	Create(context.Context, Chat) (primitive.ObjectID, error)
	FetchById(context.Context, primitive.ObjectID) ([]bson.M, error)
	FindByFilter(context.Context, interface{}) error
	FetchByFilter(context.Context, primitive.D) ([]bson.M, error)
	FetchWithLatestMessage(context.Context, primitive.D) ([]bson.M, error)
	DeleteById(context.Context, primitive.ObjectID) error
	UpdateById(context.Context, primitive.D, primitive.D) error
}

type CreateChatUseCase interface {
	Create(context.Context, Chat) (primitive.ObjectID, error)
	FetchById(context.Context, primitive.ObjectID) ([]bson.M, error)
	FindByFilter(context.Context, interface{}) error
	FetchByFilter(context.Context, primitive.D) ([]bson.M, error)
}

type UserChatsUseCase interface {
	FetchWithLatestMessage(context.Context, primitive.D) ([]bson.M, error)
}

type DeleteChatUseCase interface {
	DeleteById(context.Context, primitive.ObjectID) error
}

type GroupChatUseCase interface {
	Create(context.Context, Chat) (primitive.ObjectID, error)
	FetchById(context.Context, primitive.ObjectID) ([]bson.M, error)
}

type RenameGroupChatUseCase interface {
	UpdateById(context.Context, primitive.D, primitive.D) error
}
