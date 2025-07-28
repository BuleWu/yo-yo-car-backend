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
func (c *User) GetUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	user, err := c.userApplication.GetUserById(userId)

	if err != nil {
		c.returnJSON(ctx, utils.NewHttpError("unable to get user"), http.StatusInternalServerError)
		return
	}

	c.returnJSON(ctx, user, http.StatusOK)
}

func (c *User) CreateUser(ctx *gin.Context) {
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
	c.returnJSON(ctx, data, http.StatusCreated)
}

func (c *User) UpdateUser(ctx *gin.Context) {
	var request applications.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusInternalServerError)
		return
	}

	request.UserID = ctx.Param("id")

	user, appErr := c.userApplication.UpdateUser(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}
	c.returnJSON(ctx, user, http.StatusOK)
}

func (c *User) DeleteUser(ctx *gin.Context) {
	appErr := c.userApplication.DeleteUser(ctx.Param("id"))

	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, nil, http.StatusNoContent)
}

func (c *User) UploadProfilePicture(ctx *gin.Context) {
	userID := ctx.Param("id")

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		c.returnJSON(ctx, utils.NewHttpError("invalid file upload"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	url, appErr := c.userApplication.UploadProfilePicture(userID, file, header)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, url, http.StatusOK)
}

func (c *User) ChangePassword(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		c.returnJSON(ctx, utils.NewHttpError("Unauthorized"), http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(string)

	var request applications.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusInternalServerError)
		return
	}

	appErr := c.userApplication.ChangePassword(userID, &request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, nil, http.StatusOK)
}
