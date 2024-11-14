package handlers

import (
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/services/state_machine"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// RegisterStateMachineHandlers registers the state machine handlers
// @BasePath /
// @tag.name StateMachine
func RegisterStateMachineHandlers(engine *gin.RouterGroup, service *state_machine.StateMachineService) {
	engine.POST("/state", InitStateHandler(service))
	engine.GET("/state/:chatId", getState(service))
	engine.GET("/states", listStatesHandler(service))
	engine.POST("/state/:chatId/executeTransition", executeTransitionHandler(service))

}

// InitStateHandler initializes the state machine
// @Summary Initialize state machine
// @Description Initializes a new state machine
// @Tags StateMachine
// @Accept json
// @Produce json
// @Param request body models.InitStateMachineReq true "State Machine Initialization Request"
// @Success 200 {object} models.ChatScenarioStateObj
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /state [post]
func InitStateHandler(service *state_machine.StateMachineService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req *models.InitStateMachineReq

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, "Cannot parse CreateChatScenarioStateReq form")
			return // Early exit on unmarshaling error
		}

		obj, err := service.InitStateMachine(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		ctx.JSON(http.StatusOK, obj)
	}
}

// getState retrieves the state of a chat scenario
// @Summary Get state
// @Description Retrieves the state of a chat scenario by chat ID
// @Tags StateMachine
// @Produce json
// @Param chatId path string true "Chat ID"
// @Success 200 {object} models.ChatScenarioStateObj
// @Failure 500 {string} string "Internal Server Error"
// @Router /state/{chatId} [get]
func getState(service *state_machine.StateMachineService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatId := ctx.Param("chatId")
		state, err := service.ChatScenarioStateService.GetChatScenarioState(chatId)
		if err != nil {
			log.Fatalf("Failed to get state for chatId %s: %v", chatId, err)
		}
		ctx.JSON(http.StatusOK, state)
	}
}

// listStatesHandler lists all chat scenario states
// @Summary List states
// @Description Lists all chat scenario states
// @Tags StateMachine
// @Produce json
// @Success 200 {array} models.ChatScenarioStateObj
// @Failure 500 {string} string "Internal Server Error"
// @Router /states [get]
func listStatesHandler(service *state_machine.StateMachineService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		states, err := service.ChatScenarioStateService.ListChatScenarioStates()
		if err != nil {
			log.Fatalf("Failed to get states: %v", err)
		}
		ctx.JSON(http.StatusOK, states)
	}
}

// executeTransitionHandler executes a transition for a chat scenario
// @Summary Execute transition
// @Description Executes a transition for a chat scenario by chat ID
// @Tags StateMachine
// @Accept json
// @Produce json
// @Param chatId path string true "Chat ID"
// @Param request body models.ExecuteTransitionReq true "Transition Request"
// @Success 200 {object} models.ChatScenarioStateObj
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /state/{chatId}/executeTransition [post]
func executeTransitionHandler(service *state_machine.StateMachineService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatId := ctx.Param("chatId")
		var executeTransitionReqV2 *models.ExecuteTransitionReq

		if err := ctx.ShouldBindJSON(&executeTransitionReqV2); err != nil {
			ctx.JSON(http.StatusBadRequest, "Cannot parse execute transition request form")
			return // Early exit on unmarshaling error
		}
		chatContextObj, err := service.ChatContextLLMService.ChatContextService.GetChatContext(chatId)

		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		state, err := service.ExecuteTransition(chatId, chatContextObj.ScenarioName, executeTransitionReqV2)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			//log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		ctx.JSON(http.StatusOK, state)
	}
}
