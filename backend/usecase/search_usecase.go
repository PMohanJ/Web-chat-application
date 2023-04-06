package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type searchUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewSearchUseCase(userRepository domain.UserRepository, timeout time.Duration) domain.SearchUseCase {
	return &searchUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (su *searchUseCase) SearchUsers(c context.Context, query string) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	return su.userRepository.FetchUsers(ctx, query)
}
