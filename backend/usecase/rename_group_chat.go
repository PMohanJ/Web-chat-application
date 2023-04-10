package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type renameGroupChatUsecase struct {
	chatRepository domain.ChatRepository
	contextTimeout time.Duration
}

func NewRenameGroupChatUseCase(chatRepository domain.ChatRepository, timeout time.Duration) domain.RenameGroupChatUseCase {
	return &renameGroupChatUsecase{
		chatRepository: chatRepository,
		contextTimeout: timeout,
	}
}

func (rgc *renameGroupChatUsecase) UpdateById(c context.Context, filter primitive.D, update primitive.D) error {
	ctx, cancel := context.WithTimeout(c, rgc.contextTimeout)
	defer cancel()

	return rgc.chatRepository.UpdateById(ctx, filter, update)
}
