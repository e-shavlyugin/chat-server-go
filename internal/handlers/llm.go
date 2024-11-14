package handlers

import (
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/services/chat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterLLMHandlers(engine *gin.RouterGroup, llmService *chat.LLMService) {
	engine.POST("/llm", promptLLMHandler(llmService))
}

func promptLLMHandler(llmService *chat.LLMService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var llmPromptReq models.LLMPromptReq

		if err := ctx.ShouldBindJSON(&llmPromptReq); err != nil {
			ctx.JSON(http.StatusBadRequest, "Cannot parse llmPromptReq request form")
			return // Early exit on unmarshaling error
		}

		completion, err := llmService.GenerateFromPrompt(llmPromptReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, completion)
	}
}
