package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type getMessagesUseCase struct {
	messageRepository domain.MessageRepository
	contextTimeout    time.Duration
}

func NewGetMessagesUseCase(messageRepository domain.MessageRepository, timeout time.Duration) domain.GetMessagesUseCase {
	return &getMessagesUseCase{
		messageRepository: messageRepository,
		contextTimeout:    timeout,
	}
}

func (gm *getMessagesUseCase) FetchById(c context.Context, id primitive.ObjectID) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, gm.contextTimeout)
	defer cancel()

	return gm.messageRepository.FetchById(ctx, id)
}
