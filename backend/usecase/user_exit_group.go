package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userExitGroupUseCase struct {
	chatRepository domain.ChatRepository
	contextTimeout time.Duration
}

func NewUserExitGroupUseCase(chatRepository domain.ChatRepository, timeout time.Duration) domain.UserExitGroupUseCase {
	return &userExitGroupUseCase{
		chatRepository: chatRepository,
		contextTimeout: timeout,
	}
}

func (ueg *userExitGroupUseCase) FetchById(c context.Context, id primitive.ObjectID) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, ueg.contextTimeout)
	defer cancel()

	return ueg.chatRepository.FetchById(ctx, id)
}

func (ueg *userExitGroupUseCase) DeleteById(c context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(c, ueg.contextTimeout)
	defer cancel()

	return ueg.chatRepository.DeleteById(ctx, id)
}

func (ueg *userExitGroupUseCase) UpdateById(c context.Context, filter primitive.D, update primitive.D) error {
	ctx, cancel := context.WithTimeout(c, ueg.contextTimeout)
	defer cancel()

	return ueg.chatRepository.UpdateById(ctx, filter, update)
}
