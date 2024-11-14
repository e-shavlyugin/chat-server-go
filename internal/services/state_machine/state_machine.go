package state_machine

import (
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/services/chat"
	"fmt"
)

type StateMachineService struct {
	ChatContextLLMService    *chat.ChatContextLLMService
	ChatScenarioService      *chat.ChatScenarioService
	ChatScenarioStateService *chat.ChatScenarioStateService
	ExecutableScenarios      map[string]ExecutableScenario
	//NsbGroupStateMachine     *NsbGroupStateMachineService
}

type ExecutableScenario interface {
	GetName() string
	ExecuteScenarioStateTransition(chatId string, req *models.ExecuteTransitionReq) (*models.State, error)
}

type ExecutableScenarioRole interface {
	ExecuteScenarioRoleStateTransition(chatID string, req *models.ExecuteTransitionReq) (*models.State, error)
}

func NewStateMachine(chatContextLLMService *chat.ChatContextLLMService, chatScenarioService *chat.ChatScenarioService, chatScenarioStateService *chat.ChatScenarioStateService) (*StateMachineService, error) {
	executableScenariosMap := make(map[string]ExecutableScenario)
	//for _, executableScenario := range executableScenarios {
	//	executableScenariosMap[executableScenario.GetName()] = executableScenario
	//}
	return &StateMachineService{ChatContextLLMService: chatContextLLMService,
		ChatScenarioService:      chatScenarioService,
		ChatScenarioStateService: chatScenarioStateService,
		ExecutableScenarios:      executableScenariosMap}, nil
}

func (service *StateMachineService) RegisterScenarioStateMachineService(executableScenario ExecutableScenario) {
	service.ExecutableScenarios[executableScenario.GetName()] = executableScenario
}

func (service *StateMachineService) InitStateMachine(req *models.InitStateMachineReq) (*models.ChatScenarioStateObj, error) {
	chatContextObject, err := service.ChatContextLLMService.ChatContextService.CreateChatContext(&models.ChatContextObj{
		Name:         req.Name,
		Description:  req.Description,
		ScenarioName: req.ScenarioName,
	})
	if err != nil {
		return nil, err
	}

	state, err := service.ExecuteTransition(chatContextObject.ID, chatContextObject.ScenarioName, &models.ExecuteTransitionReq{RoleAction: models.RoleAction{ActiveRoleName: "evi", ActionName: "initialize"}})
	if err != nil {
		return nil, err
	}

	chatScenarioState, err := service.ChatScenarioStateService.CreateChatScenarioState(&models.ChatScenarioStateObj{ChatId: chatContextObject.ID, State: *state})

	if err != nil {
		return nil, err
	}
	return chatScenarioState, nil
}

func (service *StateMachineService) getChatScenarioState(chatId string) (*models.ChatScenarioStateObj, error) {
	currentState, err := service.ChatScenarioStateService.GetChatScenarioState(chatId)
	if err != nil {
		return nil, err
	}
	return currentState, nil
}

func (service *StateMachineService) ExecuteTransition(chatId string, scenarioName string, req *models.ExecuteTransitionReq) (*models.State, error) {
	currentState, err := service.getChatScenarioState(chatId)
	if err != nil {
		return nil, err
	}
	nextState, err := service.ExecutableScenarios[scenarioName].ExecuteScenarioStateTransition(chatId, req)
	if err != nil {
		if currentState == nil {
			return nil, fmt.Errorf("error executing state transition for chatId: %v and scenarioName: %v. Request: %v. Previous state was not found as well, error: %w", chatId, scenarioName, req, err)

		}
		return &currentState.State, err
	}

	//TODO think of optimising distributed transaction
	if nextState == nil {
		return nil, fmt.Errorf("error executing state transition for chatId: %v and scenarioName: %v. Request: %v. Previous state was not found as well, error: %w", chatId, scenarioName, req, err)
	}
	_, err = service.ChatScenarioStateService.UpdateChatScenarioState(&models.ChatScenarioStateObj{ChatId: chatId, State: *nextState})
	if err != nil {
		return &currentState.State, err
	}

	return nextState, nil
}
