package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type groupChatUseCase struct {
	chatRepository domain.ChatRepository
	contextTimeout time.Duration
}

func NewGroupChatUseCase(chatRepository domain.ChatRepository, timeout time.Duration) domain.GroupChatUseCase {
	return &groupChatUseCase{
		chatRepository: chatRepository,
		contextTimeout: timeout,
	}
}

func (cc *groupChatUseCase) Create(c context.Context, chat domain.Chat) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, cc.contextTimeout)
	defer cancel()

	return cc.chatRepository.Create(ctx, chat)
}

func (cc *groupChatUseCase) FetchById(c context.Context, id primitive.ObjectID) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, cc.contextTimeout)
	defer cancel()

	return cc.chatRepository.FetchById(ctx, id)
}
