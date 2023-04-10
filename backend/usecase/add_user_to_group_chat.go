package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type addUserToGroupChatUseCase struct {
	chatRepository domain.ChatRepository
	contextTimeout time.Duration
}

func NewAddUserToGroupChatUseCase(chatRepository domain.ChatRepository, timeout time.Duration) domain.AddUserToGroupChatUseCase {
	return &addUserToGroupChatUseCase{
		chatRepository: chatRepository,
		contextTimeout: timeout,
	}
}

func (aug *addUserToGroupChatUseCase) UpdateById(c context.Context, filter primitive.D, update primitive.D) error {
	ctx, cancel := context.WithTimeout(c, aug.contextTimeout)
	defer cancel()

	return aug.chatRepository.UpdateById(ctx, filter, update)
}

func (aug *addUserToGroupChatUseCase) FetchById(c context.Context, id primitive.ObjectID) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, aug.contextTimeout)
	defer cancel()

	return aug.chatRepository.FetchById(ctx, id)
}
