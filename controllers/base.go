package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct{}

func (c *Controller) returnJSONError(ctx *gin.Context, message string) {
	c.returnJSON(ctx, message, http.StatusBadRequest)
}

func (c *Controller) returnJSON(ctx *gin.Context, data interface{}, status int) {
	ctx.JSON(status, data)
}
