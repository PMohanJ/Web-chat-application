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

type LoginController struct {
	LoginUseCase domain.LoginUseCase
	Env          *bootstrap.Env
}

func (lc *LoginController) Login(c *gin.Context) {
	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while decoding user data"})
		log.Println(err)
		return
	}

	var registeredUser domain.User

	// check if user is a registered user
	registeredUser, err := lc.LoginUseCase.GetByEmail(c, user.Email)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not registered"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while querying user data"})
		log.Panic(err)
	}

	// user exist, check for password validation
	resgisteredPassword := registeredUser.Password
	errMsg, valid := helpers.VerifyPassword(resgisteredPassword, user.Password)
	if !valid {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: errMsg})
		return
	}

	id := registeredUser.Id.Hex()
	if registeredUser.Token, err = lc.LoginUseCase.CreateAccessToken(id, user.Name, user.Email, lc.Env.SecretKey); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Unable to generate token"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, registeredUser)
}
