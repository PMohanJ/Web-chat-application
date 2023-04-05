package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type singupUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewSignupUseCase(userRepository domain.UserRepository, timeout time.Duration) domain.SignupUsecase {
	return &singupUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (su *singupUseCase) Create(c context.Context, user domain.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	return su.userRepository.Create(ctx, user)
}

func (su *singupUseCase) GetByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	return su.userRepository.GetByEmail(ctx, email)
}

func (su *singupUseCase) CreateAccessToken(id, name, email, secret string) (string, error) {
	return helpers.GenerateToken(id, name, email, secret)
}
