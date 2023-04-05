package domain

import "context"

type LoginUseCase interface {
	GetByEmail(context.Context, string) (User, error)
	CreateAccessToken(string, string, string, string) (string, error)
}
