package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetQueryInt(ctx *gin.Context, param string, defaultVal int) int {
	paramStr := ctx.Query(param) // read query param
	if paramStr == "" {
		return defaultVal
	}

	paramVal, err := strconv.Atoi(paramStr)
	if err != nil {
		return defaultVal
	}

	return paramVal
}
