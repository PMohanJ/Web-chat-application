package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type editMessageUseCase struct {
	messageRepository domain.MessageRepository
	contextTimeout    time.Duration
}

func NewEditMessageUseCase(messageRepository domain.MessageRepository, timeout time.Duration) domain.EditMessageUseCase {
	return &editMessageUseCase{
		messageRepository: messageRepository,
		contextTimeout:    timeout,
	}
}

func (em *editMessageUseCase) UpdateById(c context.Context, filter primitive.D, update primitive.D) error {
	ctx, cancel := context.WithTimeout(c, em.contextTimeout)
	defer cancel()

	return em.messageRepository.UpdateById(ctx, filter, update)
}

func (em *editMessageUseCase) FetchById(c context.Context, field string, id primitive.ObjectID) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, em.contextTimeout)
	defer cancel()

	return em.messageRepository.FetchById(ctx, field, id)
}
