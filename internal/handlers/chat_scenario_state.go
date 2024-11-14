package handlers

import (
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/services/chat"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterChatScenarioStateHandlers(engine *gin.RouterGroup, dialogueScenarioService *chat.ChatScenarioStateService) {
	engine.POST("/chatScenarioState", createChatScenarioStateHandler(dialogueScenarioService))
	engine.GET("/chatScenarioState/:chatId", getChatScenarioStateHandler(dialogueScenarioService))
	engine.GET("/chatScenarioStates", listChatScenarioStateHandler(dialogueScenarioService))
}

func createChatScenarioStateHandler(service *chat.ChatScenarioStateService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req *models.ChatScenarioStateObj

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, "Cannot parse CreateChatScenarioStateReq form")
			return // Early exit on unmarshaling error
		}

		chatScenarioStateObject, err := service.CreateChatScenarioState(req)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		ctx.JSON(http.StatusOK, chatScenarioStateObject)
	}
}

func getChatScenarioStateHandler(service *chat.ChatScenarioStateService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatId := ctx.Param("chatId")
		state, err := service.GetChatScenarioState(chatId)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		ctx.JSON(http.StatusOK, state)
	}
}

func listChatScenarioStateHandler(service *chat.ChatScenarioStateService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		states, err := service.ListChatScenarioStates()
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		ctx.JSON(http.StatusOK, states)
	}
}
