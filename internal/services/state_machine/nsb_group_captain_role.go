package state_machine

import (
	"chat_server_v2/internal/models"
	"fmt"
)

type NsbGroupCaptainRoleStateMachineService struct {
	Role                *models.Role
	StateMachineService *StateMachineService
}

func NewNsbGroupCaptainRoleStateMachineService(captainRole *models.Role, service *StateMachineService) *NsbGroupCaptainRoleStateMachineService {
	return &NsbGroupCaptainRoleStateMachineService{Role: captainRole, StateMachineService: service}
}

func (service *NsbGroupCaptainRoleStateMachineService) ExecuteScenarioRoleStateTransition(chatId string, req *models.ExecuteTransitionReq) (*models.State, error) {
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

	switch req.RoleAction.ActionName {
	case "send_message":
		nextState, _ = service.sendMessage(chatId, action, &req.Payload)
	default:
		return nil, fmt.Errorf("action '%s' not found for role '%s'", req.RoleAction.ActionName, req.RoleAction.ActiveRoleName)
	}
	return nextState, nil
}

func (service *NsbGroupCaptainRoleStateMachineService) sendMessage(chatId string, action *models.Action, payload *models.TransitionRequestPayload) (*models.State, error) {
	nextState := models.State{NextActiveRoleName: action.NextRole, NextActions: action.NextGuaranteedActions}
	chatContextObject, err := service.StateMachineService.ChatContextLLMService.ChatContextService.GetChatContext(chatId)
	if err != nil {
		return nil, err
	}
	prompt, err := service.StateMachineService.ChatContextLLMService.SendPrompt(chatContextObject, service.Role.Name, action.LLMConfig.SystemPromptTemplate, payload.Message)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Role: %s\n. Prompt: %s\n", service.Role, prompt)
	return &nextState, nil
}
