package chatControllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/bootstrap"
	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/helpers"
	"go.mongodb.org/mongo-driver/mongo"
)

type SingUpController struct {
	SingupUseCase domain.SignupUsecase
	Env           *bootstrap.Env
}

func (sc *SingUpController) SingUp(c *gin.Context) {
	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while decoding user data"})
		log.Println(err)
		return
	}

	if user.Pic == "" {
		user.SetDefaultPic()
	}

	var temp domain.User

	temp, err := sc.SingupUseCase.GetByEmail(c, user.Email)

	// if err is other than ErrNoDocuments, something wrong while querying
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while querying for user"})
		log.Panic(err)
	}

	if temp.Email == user.Email {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "You've already registered with this email"})
		return
	}

	// user doesn't exist in database, so register the user
	user.Password = helpers.HashPassowrd(user.Password)

	insertedId, err := sc.SingupUseCase.Create(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while registering the user"})
		log.Panic(err)
	}

	user.Id = insertedId

	id := user.Id.Hex()
	if user.Token, err = sc.SingupUseCase.CreateAccessToken(id, user.Name, user.Email, sc.Env.SecretKey); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Unable to generate token"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, user)
}
