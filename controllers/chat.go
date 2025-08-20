package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/core/utils"
)

func NewChatController(
	chatApplication *applications.Chat,
) *Chat {
	return &Chat{
		chatApplication: chatApplication,
	}
}

type Chat struct {
	Controller
	chatApplication *applications.Chat
}

// GetUserChats GET /chats?userId=:id
func (c *Chat) GetUserChats(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	chats, appErr := c.chatApplication.GetUserChats(userID)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, chats, http.StatusOK)
}

// GetChat GET /chats/:id
func (c *Chat) GetChat(ctx *gin.Context) {
	chat, appErr := c.chatApplication.GetChat(ctx.Param("id"))

	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, chat, http.StatusOK)
}

// CreateChat POST /chats
func (c *Chat) CreateChat(ctx *gin.Context) {
	var request applications.CreateChatRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusBadRequest)
		return
	}

	chat, status, appErr := c.chatApplication.CreateChat(&request)
	if appErr != nil {
		c.returnJSON(ctx, appErr, status)
		return
	}

	c.returnJSON(ctx, chat, status)
}

// DeleteChat DELETE /chats/:id
func (c *Chat) DeleteChat(ctx *gin.Context) {
	chatID := ctx.Param("id")

	if appErr := c.chatApplication.DeleteChat(chatID); appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, nil, http.StatusNoContent)
}

// GetChatMessages GET /chats/:chatId/messages
func (c *Chat) GetChatMessages(ctx *gin.Context) {
	chatID := ctx.Param("chatId")

	limit := utils.GetQueryInt(ctx, "limit", 20)
	offset := utils.GetQueryInt(ctx, "offset", 0)

	messages, appErr := c.chatApplication.GetChatMessages(chatID, limit, offset)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, messages, http.StatusOK)
}

// SendMessage POST /chats/:chatId/messages
func (c *Chat) SendMessage(ctx *gin.Context) {
	chatID := ctx.Param("id")

	var request applications.SendMessageRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusBadRequest)
		return
	}
	request.ChatID = chatID

	message, appErr := c.chatApplication.SendMessage(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, message, http.StatusCreated)
}
