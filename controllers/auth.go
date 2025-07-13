package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/core/utils"
)

func NewAuthController(
	userApplication *applications.User,
) *Auth {
	return &Auth{
		userApplication: userApplication,
	}
}

type Auth struct {
	Controller
	userApplication *applications.User
}

/*func (c Auth) Login(ctx *gin.Context) {
	var credentials struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var token string

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusBadRequest)
		return
	}

	user, err := c.userApplication.GetUserByEmail(credentials.Email)
	if err != nil {
		c.returnJSON(ctx, utils.NewHttpError("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	if !utils.VerifyPassword(credentials.Password, user.Password) {
		c.returnJSON(ctx, utils.NewHttpError("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	fmt.Println("The user email: ", credentials.Email)
	fmt.Println("The user password: ", credentials.Password)

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		c.returnJSON(ctx, , http.StatusInternalServerError)
		return
	}

	c.returnJSON(ctx, token, http.StatusOK)
}*/

func (c Auth) Register(ctx *gin.Context) {
	var credentials struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusBadRequest)
		return
	}

	user, _ := c.userApplication.GetUserByEmail(credentials.Email)

	if user != nil {
		c.returnJSON(ctx, utils.NewHttpError("a user with this email already exists"), http.StatusNotAcceptable)
		return
	}

	fmt.Println("User email: ", credentials.Email)
	fmt.Println("User pass: ", credentials.Password)

	if len(credentials.Password) < 8 {
		c.returnJSON(ctx, utils.NewHttpError("the user password needs to be at least 8 characters long"), http.StatusNotAcceptable)
		return
	}

	hashedPassword, _ := utils.HashPassword(credentials.Password)

	req := &applications.CreateUserRequest{
		FirstName: credentials.FirstName,
		LastName:  credentials.LastName,
		Email:     credentials.Email,
		Password:  hashedPassword,
	}

	user, appErr := c.userApplication.CreateUser(req)
	if appErr != nil {
		c.returnJSON(ctx, appErr.GetMessage(), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, utils.NewHttpError("registered successfully"), http.StatusOK)
}
