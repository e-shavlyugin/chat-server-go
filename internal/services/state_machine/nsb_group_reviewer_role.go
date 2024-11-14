package state_machine

import (
	"chat_server_v2/internal/models"
)

type NsbGroupReviewerRoleStateMachineService struct {
	Role                *models.Role
	StateMachineService *StateMachineService
}

func NewNsbGroupReviewerRoleStateMachineService(eviRole *models.Role, service *StateMachineService) *NsbGroupReviewerRoleStateMachineService {
	return &NsbGroupReviewerRoleStateMachineService{Role: eviRole, StateMachineService: service}
}

func (stateMachine *NsbGroupReviewerRoleStateMachineService) ExecuteScenarioRoleStateTransition(chatId string, req *models.ExecuteTransitionReq) (*models.State, error) {
	return nil, nil
}

/*
func (stateMachine *NsbGroupReviewerRoleStateMachineService) initialize(req *models.ExecuteTransitionReq) (*models.State, error) {
	nextState := models.State{NextActiveRoleName: "captain"}
	nextState.NextActions = []string{"send_message"}

	chatContextObject, err := stateMachine.ChatScenarioStateService.ChatContextLLMService.ChatContextService.GetChatContext(req.ChatId)
	if err != nil {
		return nil, err
	}
	stateMachine.ChatScenarioStateService.ChatContextLLMService.SendAIMessage(chatContextObject, req.RoleAction.ActiveRoleName, stateMachine.RoleDefinition.Actions["initialise"].LLMConfig.SystemPromptTemplate, "Hello captain, how can I help yoh today?")

	return &nextState, nil
}

*/
