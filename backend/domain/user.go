package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionUser = "user"

type User struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"password" bson:"password"`
	Pic        string             `json:"pic" bson:"pic"`
	IsAdmin    bool               `json:"isAdmin" bson:"isAdmin"`
	Created_at time.Time          `json:"-" bson:"created_at"`
	Updated_at time.Time          `json:"-" bson:"updated_at"`
	Token      string             `json:"token" bson:"-"`
}

type UserRepository interface {
	Create(context.Context, User) (primitive.ObjectID, error)
	GetByEmail(context.Context, string) (User, error)
	FetchUsers(context.Context, string) ([]bson.M, error)
}

func (u *User) SetDefaultPic() {
	u.Pic = "https://res.cloudinary.com/dkqc4za4f/image/upload/v1671523788/default_toic85.png"
}
