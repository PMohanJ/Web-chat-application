package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chat struct {
	Id          primitive.ObjectID   `bson:"_id"`
	IsGroupChat bool                 `bson:"isGroupChat"`
	Users       []primitive.ObjectID `bson:"users"`
	ChatName    string               `bson:"chatName"`
}
