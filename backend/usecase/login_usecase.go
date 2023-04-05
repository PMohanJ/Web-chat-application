package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/helpers"
)

type loginUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLoginUseCase(userRepository domain.UserRepository, timeout time.Duration) domain.LoginUseCase {
	return &loginUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (lu *loginUseCase) GetByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()

	return lu.userRepository.GetByEmail(ctx, email)
}

func (lu *loginUseCase) CreateAccessToken(id, name, email, secret string) (string, error) {
	return helpers.GenerateToken(id, name, email, secret)
}
