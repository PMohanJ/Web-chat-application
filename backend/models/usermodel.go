package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,required" bson:"name"`
	Email    string             `jaon:"email,required" bson:"email"`
	Password string             `json:"password,required" bson:"password"`
	Pic      string             `json:"pic" bson:"pic"`
	IsAdmin  bool               `json:"isAdmin" bson:"isAdmin"` // if no value is provided, then bydefault it is set to false
}

func (u *User) SetDefaultPic() {
	u.Pic = "https://p.kindpng.com/picc/s/24-248253_user-profile-default-image-png-clipart-png-download.png"
}
