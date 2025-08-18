package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/core/utils"
)

func NewMessageController(
	messageApplication *applications.Message,
) *Message {
	return &Message{
		messageApplication: messageApplication,
	}
}

type Message struct {
	Controller
	messageApplication *applications.Message
}

func (c *Message) UpdateMessage(ctx *gin.Context) {
	var request applications.UpdateMessageRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusInternalServerError)
		return
	}

	request.MessageID = ctx.Param("id")

	message, appErr := c.messageApplication.UpdateMessage(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, message, http.StatusOK)
}
