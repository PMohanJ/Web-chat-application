package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type sendMessageUseCase struct {
	messageRepository domain.MessageRepository
	chatRepository    domain.ChatRepository
	contextTimeout    time.Duration
}

func NewSendMessageUseCase(chatRepository domain.ChatRepository, messageRepository domain.MessageRepository, timeout time.Duration) domain.SendMessageUseCase {
	return &sendMessageUseCase{
		messageRepository: messageRepository,
		chatRepository:    chatRepository,
		contextTimeout:    timeout,
	}
}

func (sm *sendMessageUseCase) Create(c context.Context, message domain.Message) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, sm.contextTimeout)
	defer cancel()

	return sm.messageRepository.Create(ctx, message)
}

func (sm *sendMessageUseCase) UpdateByFilter(c context.Context, filter primitive.D, update primitive.D) error {
	ctx, cancel := context.WithTimeout(c, sm.contextTimeout)
	defer cancel()

	return sm.chatRepository.UpdateById(ctx, filter, update)
}

func (sm *sendMessageUseCase) FetchById(c context.Context, field string, id primitive.ObjectID) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, sm.contextTimeout)
	defer cancel()

	return sm.messageRepository.FetchById(ctx, field, id)
}
