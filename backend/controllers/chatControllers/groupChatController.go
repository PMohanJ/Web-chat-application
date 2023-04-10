package chatControllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupChatController struct {
	GroupChatUseCase domain.GroupChatUseCase
}

func (gc *GroupChatController) CreateGroupChat(c *gin.Context) {
	var groupData map[string]interface{}

	if err := c.BindJSON(&groupData); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while parsing data"})
		log.Panic(err)
	}

	groupName := groupData["groupName"].(string)

	groupPic, ok := groupData["groupPic"].(string)
	if !ok || groupPic == "" {
		groupPic = domain.GetDefaultGroupPic()
	}

	users := groupData["users"].([]interface{})

	aUser, exists := c.Get("_id")
	if !exists {
		log.Panic("User details not available")
	}
	adminUser := aUser.(primitive.ObjectID)

	var usersIds []primitive.ObjectID
	usersIds = append(usersIds, adminUser)

	for _, uId := range users {
		id := uId.(string)

		temp, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Panic(err)
		}
		usersIds = append(usersIds, temp)
	}

	groupChat := domain.Chat{
		IsGroupChat: true,
		ChatName:    groupName,
		Users:       usersIds,
		GroupAdmin:  adminUser,
		GroupPic:    groupPic,
	}

	insId, err := gc.GroupChatUseCase.Create(c, groupChat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while inserting data"})
		log.Panic(err)
	}

	insertedChat, err := gc.GroupChatUseCase.FetchById(c, insId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while retrieving data"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, insertedChat[0])
}
