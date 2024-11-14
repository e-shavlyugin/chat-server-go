package handlers

import (
	"chat_server_v2/internal/services/chat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterChatScenarioHandlers(engine *gin.RouterGroup, chatScenarioService *chat.ChatScenarioService) {
	engine.GET("/scenario/:scenarioName", getChatScenarioHandler(chatScenarioService))
	engine.GET("/scenarios", listChatScenariosHandler(chatScenarioService))
}

func getChatScenarioHandler(chatScenarioService *chat.ChatScenarioService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		scenarioName := ctx.Param("scenarioName")
		scenarioObject, err := chatScenarioService.GetScenario(scenarioName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, scenarioObject)
	}
}

func listChatScenariosHandler(chatScenarioService *chat.ChatScenarioService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roles, err := chatScenarioService.ListScenarios()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, roles)
	}
}
