package controllers

import (
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

func (c Auth) Login(ctx *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var token string

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusBadRequest)
		return
	}

	user, err := c.userApplication.GetUserByEmail(credentials.Email)
	if err != nil || user == nil {
		c.returnJSON(ctx, utils.NewHttpError("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	if !utils.VerifyPassword(user.Password, credentials.Password) {
		c.returnJSON(ctx, utils.NewHttpError("invalid credentials"), http.StatusUnauthorized)
		return
	}

	token, appErr := utils.GenerateToken(user.ID)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError("error generating token"), http.StatusInternalServerError)
		return
	}

	c.returnJSON(ctx, token, http.StatusOK)
}

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
