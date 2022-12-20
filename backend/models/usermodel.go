package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,required" bson:"name"`
	Email      string             `json:"email,required" bson:"email"`
	Password   string             `json:"password,required" bson:"password"`
	Pic        string             `json:"pic" bson:"pic"`
	IsAdmin    bool               `json:"isAdmin" bson:"isAdmin"` // if no value is provided, then bydefault it is set to false
	Created_at time.Time          `json:"-" bson:"created_at"`
	Updated_at time.Time          `json:"-" bson:"updated_at"`
	Token      string             `json:"token" bson:"-"`
}

type UserResponse struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"-" bson:"password"`
	Pic      string             `json:"pic" bson:"pic"`
	IsAdmin  bool               `json:"-" bson:"isAdmin"`
	Token    string             `json:"token" bson:"-"`
}

func (u *User) SetDefaultPic() {
	u.Pic = "https://res.cloudinary.com/dkqc4za4f/image/upload/v1671523788/default_toic85.png"
}
