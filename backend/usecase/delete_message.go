package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type deleteMessageUseCase struct {
	messageRepository domain.MessageRepository
	contextTimeout    time.Duration
}

func NewDeleteMessageUseCase(messageRepository domain.MessageRepository, timeout time.Duration) domain.DeleteMessageUseCase {
	return &deleteMessageUseCase{
		messageRepository: messageRepository,
		contextTimeout:    timeout,
	}
}

func (dm *deleteMessageUseCase) DeleteById(c context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(c, dm.contextTimeout)
	defer cancel()

	return dm.messageRepository.DeleteById(ctx, id)
}
