package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/core/utils"
	"zavrsni/yo-yo-car/runtimebag"
)

var googleOAuthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     runtimebag.GetEnvString("GOOGLE_CLIENT_ID", ""),
	ClientSecret: runtimebag.GetEnvString("GOOGLE_CLIENT_SECRET", ""),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

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

func (c Auth) OauthGoogleLogin(ctx *gin.Context) {
	oauthState, err := generateStateOAuthCookie()
	if err != nil {
		c.returnJSON(ctx, utils.NewHttpError("unable to generate state oauth cookie"), http.StatusInternalServerError)
	}

	u := googleOAuthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOnline, oauth2.ApprovalForce)
	ctx.Redirect(http.StatusTemporaryRedirect, u)
}

type TestRequest struct {
	Code string `json:"code"`
}

func (c Auth) OauthGoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")

	data, err := getUserDataFromGoogle(code)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/welcome")
		return
	}

	fmt.Println("User data", string(data))

	type googleUserData struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}

	var userData googleUserData

	if err = json.Unmarshal(data, &userData); err != nil {
		c.returnJSON(ctx, utils.NewHttpError("unable to unmarshal the user data"), http.StatusInternalServerError)
		return
	}

	user, _ := c.userApplication.GetUserByEmail(userData.Email)

	var request = &applications.CreateUserRequest{
		FirstName: userData.FamilyName,
		LastName:  userData.GivenName,
		Email:     userData.Email,
		Password:  "",
	}

	if user == nil {
		c.userApplication.CreateUser(request)
		return
	}
}

func generateStateOAuthCookie() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)

	return state, nil
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange went wrong: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %s", err.Error())
	}

	return contents, nil
}
