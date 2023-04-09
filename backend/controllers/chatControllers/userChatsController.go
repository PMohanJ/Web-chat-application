package chatControllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserChatsController struct {
	UserChatsUseCase domain.UserChatsUseCase
}

func (uc *UserChatsController) GetUserChats(c *gin.Context) {
	id, exists := c.Get("_id")
	if !exists {
		log.Panic("User details not available")
	}

	userId := id.(primitive.ObjectID)

	matchStage := bson.D{
		{
			"$match", bson.D{
				{
					"users", bson.D{{"$elemMatch", bson.D{{"$eq", userId}}}},
				},
			},
		},
	}

	userChats, err := uc.UserChatsUseCase.FetchWithLatestMessage(c, matchStage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error in the server"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, userChats)
}
