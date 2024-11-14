package middleware

import (
	"chat_server_v2/internal/models"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
}

func AbortWithErrorObject(ctx *gin.Context, status int, err error) {
	errJson := models.Errors{
		Errors: []string{err.Error()},
	}
	ctx.AbortWithStatusJSON(status, errJson)
}
