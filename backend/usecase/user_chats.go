package usecase

import (
	"context"
	"time"

	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userChatsUseCase struct {
	chatRepository domain.ChatRepository
	contexTimeout  time.Duration
}

func NewUserChatsUseCase(chatRepository domain.ChatRepository, timeout time.Duration) domain.UserChatsUseCase {
	return &userChatsUseCase{
		chatRepository: chatRepository,
		contexTimeout:  timeout,
	}
}

func (uc *userChatsUseCase) FetchWithLatestMessage(c context.Context, filter primitive.D) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(c, uc.contexTimeout)
	defer cancel()

	return uc.chatRepository.FetchWithLatestMessage(ctx, filter)
}
