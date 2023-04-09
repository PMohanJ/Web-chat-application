package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createChatUseCase struct {
	chatRepository domain.ChatRepository
	contextTimeout time.Duration
}

func NewCreateChatUseCase(chatRepository domain.ChatRepository, timeout time.Duration) domain.CreateChatUseCase {
	return &createChatUseCase{
		chatRepository: chatRepository,
		contextTimeout: timeout,
	}
}

func (cc *createChatUseCase) Create(c context.Context, chat domain.Chat) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, cc.contextTimeout)
	defer cancel()

	return cc.chatRepository.Create(ctx, chat)
}

func (cc *createChatUseCase) FetchById(c context.Context, id primitive.ObjectID) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, cc.contextTimeout)
	defer cancel()

	return cc.chatRepository.FetchById(ctx, id)
}

func (cc *createChatUseCase) FindByFilter(c context.Context, filter interface{}) error {
	ctx, cancel := context.WithTimeout(c, cc.contextTimeout)
	defer cancel()

	return cc.chatRepository.FindByFilter(ctx, filter)
}

func (cc *createChatUseCase) FetchByFilter(c context.Context, filter primitive.D) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, cc.contextTimeout)
	defer cancel()

	return cc.chatRepository.FetchByFilter(ctx, filter)
}
