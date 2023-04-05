package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignupUsecase interface {
	Create(context.Context, User) (primitive.ObjectID, error)
	GetByEmail(context.Context, string) (User, error)
	CreateAccessToken(string, string, string, string) (string, error)
}
