package server

import (
	"chat_server_v2/config"
	"chat_server_v2/internal/handlers"
	"chat_server_v2/internal/repositories"
	"chat_server_v2/internal/services/chat"
	"chat_server_v2/internal/services/state_machine"
	"chat_server_v2/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine

	chatLLMService           *chat.ChatContextLLMService
	chatContextService       *chat.ChatContextService
	llmService               *chat.LLMService
	chatScenarioService      *chat.ChatScenarioService
	chatScenarioStateService *chat.ChatScenarioStateService
	stateMachineService      *state_machine.StateMachineService
}

func NewServer(cfg *config.Config) (*Server, error) {

	var chatContextRepository = repositories.NewChatContextsRepository(cfg)
	chatContextService := chat.NewChatContextService(chatContextRepository)
	var llmService = chat.NewLLMService("openai")
	var chatContextLLMService = chat.NewChatContextLLMService(chatContextService, llmService)

	var chatScenarioRepository, err1 = repositories.NewChatScenarioRepository("file")
	if err1 != nil {
		return nil, err1
	}
	chatScenarioService := chat.NewChatScenarioService(chatScenarioRepository)

	chatScenarioStateRepository, _ := repositories.NewChatScenarioStateRepository(cfg)
	chatScenarioStateService := chat.NewChatScenarioStateService(chatScenarioStateRepository, chatContextService)

	//TODO check logic
	stateMachineService, _ := state_machine.NewStateMachine(chatContextLLMService, chatScenarioService, chatScenarioStateService)
	//tmp, _ := stateMachineService.ChatScenarioService.GetScenario("nsbGroup")
	//fmt.Printf("tmp: %v\n", tmp)
	nsbGroupStateMachineService, _ := state_machine.NewNsbGroupStateMachineService(stateMachineService)
	stateMachineService.RegisterScenarioStateMachineService(nsbGroupStateMachineService)

	router := gin.Default()
	mainGroup := router.Group("/v1")
	server := &Server{
		router:                   router,
		chatScenarioStateService: chatScenarioStateService,
		chatLLMService:           chatContextLLMService,
		chatContextService:       chatContextService,
		llmService:               llmService,
		chatScenarioService:      chatScenarioService,
		stateMachineService:      stateMachineService,
	}

	server.router.Use(middleware.SetHeader)
	server.router.Use(middleware.CreateCorsMiddleware(cfg.HTTP.Cors))

	handlers.RegisterSwaggerHandler(mainGroup)
	handlers.RegisterChatScenarioStateHandlers(mainGroup, server.chatScenarioStateService)
	handlers.RegisterChatContextHandlers(mainGroup, server.chatContextService)
	handlers.RegisterChatScenarioHandlers(mainGroup, server.chatScenarioService)
	handlers.RegisterLLMHandlers(mainGroup, server.llmService)
	handlers.RegisterStateMachineHandlers(mainGroup, server.stateMachineService)
	return server, nil
}

func Start(cfg *config.Config) error {
	server, err := NewServer(cfg)
	if err != nil {
		return err
	}

	return server.router.Run(cfg.Port)
}
