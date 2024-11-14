package handlers

import (
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/services/chat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterChatContextHandlers(engine *gin.RouterGroup, chatContextService *chat.ChatContextService) {
	engine.POST("/chatContext", createChatContextHandler(chatContextService))
	engine.GET("/chatContext/:chatId", getChatContextHandler(chatContextService))
	engine.GET("/chatContexts", listChatContextHandler(chatContextService))
}

func createChatContextHandler(chatContextService *chat.ChatContextService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var createChatContextObjectReq models.CreateChatContextObjectReq

		if err := ctx.ShouldBindJSON(&createChatContextObjectReq); err != nil {
			ctx.JSON(http.StatusBadRequest, "Cannot parse createChatContextObjectReq request form")
			return
		}

		chatContextId, err := chatContextService.CreateChatContext(&models.ChatContextObj{
			ID:           createChatContextObjectReq.ID,
			Name:         createChatContextObjectReq.Name,
			Description:  createChatContextObjectReq.Description,
			ScenarioName: createChatContextObjectReq.ScenarioName,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, chatContextId)
	}
}

func getChatContextHandler(chatContextService *chat.ChatContextService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatId := ctx.Param("chatId")
		chatContextObject, err := chatContextService.GetChatContext(chatId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if chatContextObject == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Home context with ID " + chatId + " not found"})
			return
		}
		ctx.JSON(http.StatusOK, chatContextObject)
	}
}

func listChatContextHandler(chatContextService *chat.ChatContextService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatContextObjects, err := chatContextService.ListChatContexts()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, chatContextObjects)
	}
}
