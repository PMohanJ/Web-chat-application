package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type deleteChatUseCase struct {
	chatRepository domain.ChatRepository
	contextTimeOut time.Duration
}

func NewDeleteChatUseCase(chatRepository domain.ChatRepository, timeout time.Duration) domain.DeleteChatUseCase {
	return &deleteChatUseCase{
		chatRepository: chatRepository,
		contextTimeOut: timeout,
	}
}

func (dc *deleteChatUseCase) DeleteById(c context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(c, dc.contextTimeOut)
	defer cancel()

	return dc.chatRepository.DeleteById(ctx, id)
}
