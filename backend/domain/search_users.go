package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type SearchUseCase interface {
	SearchUsers(context.Context, string) ([]bson.M, error)
}
