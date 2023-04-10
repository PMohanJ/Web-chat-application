package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type removeUserFromGroupUseCase struct {
	chatRepository domain.ChatRepository
	contextTimeout time.Duration
}

func NewRemoveUserFromGroupUseCase(chatRepository domain.ChatRepository, timeout time.Duration) domain.RemoveUserFromGroupUseCase {
	return &removeUserFromGroupUseCase{
		chatRepository: chatRepository,
		contextTimeout: timeout,
	}
}

func (rug *removeUserFromGroupUseCase) UpdateById(c context.Context, filter primitive.D, update primitive.D) error {
	ctx, cancel := context.WithTimeout(c, rug.contextTimeout)
	defer cancel()

	return rug.chatRepository.UpdateById(ctx, filter, update)
}

func (rug *removeUserFromGroupUseCase) FetchById(c context.Context, id primitive.ObjectID) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, rug.contextTimeout)
	defer cancel()

	return rug.chatRepository.FetchById(ctx, id)
}
