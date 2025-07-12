package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/core/utils"
)

// NewUserController user Controller Constructor that handles dependency injection
func NewUserController(
	userApplication *applications.User,
) *User {
	return &User{
		userApplication: userApplication,
	}
}

type User struct {
	Controller
	userApplication *applications.User
}

/*GetUser is the UserController method that handles the GET request */
func (c User) GetUser(ctx *gin.Context) {

	userId := ctx.Param("id")
	data, err := c.userApplication.GetUser(userId)

	if err != nil {
		c.returnJSON(ctx, utils.NewHttpError("unable to get user"), http.StatusInternalServerError)
		return
	}

	c.returnJSON(ctx, data, http.StatusOK)
}

func (c User) CreateUser(ctx *gin.Context) {
	var request applications.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusInternalServerError)
		return
	}

	data, appErr := c.userApplication.CreateUser(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}
	c.returnJSON(ctx, data, http.StatusOK)
	return
}
