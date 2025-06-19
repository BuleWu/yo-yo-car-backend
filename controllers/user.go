package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zavrsni/yo-yo-car/core/utils"
	"zavrsni/yo-yo-car/repositories"
)

type User struct {
	Controller
	userRepository repositories.UserRepository
}

/*GetUser is the UserController method that handles the GET request */
func (c User) GetUser(ctx *gin.Context) {

	userId := ctx.Param("id")
	data, err := c.userRepository.Get(userId)

	if err != nil {
		c.returnJSON(ctx, utils.NewHttpError("unable to get user"), http.StatusInternalServerError)
		return
	}

	c.returnJSON(ctx, data, http.StatusOK)
}
