package state_machine

import (
	"chat_server_v2/internal/models"
	"fmt"
)

type NsbGroupEviRoleStateMachineService struct {
	Role                *models.Role
	StateMachineService *StateMachineService
}

func NewNsbGroupEviRoleStateMachineService(eviRole *models.Role, service *StateMachineService) *NsbGroupEviRoleStateMachineService {
	return &NsbGroupEviRoleStateMachineService{Role: eviRole, StateMachineService: service}
}

func (service *NsbGroupEviRoleStateMachineService) ExecuteScenarioRoleStateTransition(chatId string, req *models.ExecuteTransitionReq) (*models.State, error) {
	var nextState *models.State
	var action *models.Action
	for _, elem := range service.Role.Actions {
		if elem.Name == req.RoleAction.ActionName {
			action = &elem
		}
	}
	if action == nil {
		return nil, fmt.Errorf("action '%s' not found for role '%s'", req.RoleAction.ActionName, req.RoleAction.ActiveRoleName)
	}

	switch action.Name {
	case "initialize":
		nextState, _ = service.initialize(chatId, action, req)
	default:
		return nil, fmt.Errorf("action '%s' not found for role '%s'", req.RoleAction.ActionName, req.RoleAction.ActiveRoleName)
	}
	return nextState, nil
}

func (service *NsbGroupEviRoleStateMachineService) initialize(chatId string, action *models.Action, req *models.ExecuteTransitionReq) (*models.State, error) {
	nextState := models.State{NextActiveRoleName: action.NextRole, NextActions: action.NextGuaranteedActions}

	chatContextObject, err := service.StateMachineService.ChatContextLLMService.ChatContextService.GetChatContext(chatId)
	if err != nil {
		return nil, err
	}
	service.StateMachineService.ChatContextLLMService.SendAIMessage(chatContextObject, req.RoleAction.ActiveRoleName, action.LLMConfig.SystemPromptTemplate, "Hello captain, how can I help yoh today?")
	return &nextState, nil
}
