package chatControllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserExitGroupController struct {
	UserExitGroupUseCase domain.UserExitGroupUseCase
}

func (ueg *UserExitGroupController) UserExitGroup(c *gin.Context) {
	var reqData map[string]interface{}

	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while parsing data"})
		log.Panic(err)
	}

	cId := reqData["chatId"].(string)
	uId, exists := c.Get("_id")
	if !exists {
		log.Panic("User details not available")
	}

	chatId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		log.Panic(err)
	}

	userId := uId.(primitive.ObjectID)

	filter := bson.D{{"_id", chatId}}

	Documents, err := ueg.UserExitGroupUseCase.FetchById(c, chatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while retrieving data"})
		log.Panic(err)
	}

	chatDocu := Documents[0]
	groupAdmin := chatDocu["groupAdmin"].(primitive.ObjectID)

	// check if admin is exiting the Group
	if userId.Hex() == groupAdmin.Hex() {
		// delete the whole chat
		err := ueg.UserExitGroupUseCase.DeleteById(c, chatId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while querying database"})
			log.Panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Exited from group"})
		return
	}

	// just remove the user from Group chat
	update := bson.D{{"$pull", bson.D{{"users", userId}}}}

	err = ueg.UserExitGroupUseCase.UpdateById(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while updating data"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exited from group"})
}
