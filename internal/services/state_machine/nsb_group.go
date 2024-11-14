package state_machine

import (
	"chat_server_v2/internal/models"
	"fmt"
)

type NsbGroupStateMachineService struct {
	Scenario                        *models.ScenarioObject
	StateMachineService             *StateMachineService
	EviRoleStateMachineService      *NsbGroupEviRoleStateMachineService
	CaptainRoleStateMachineService  *NsbGroupCaptainRoleStateMachineService
	ReviewerRoleStateMachineService *NsbGroupReviewerRoleStateMachineService
}

func NewNsbGroupStateMachineService(service *StateMachineService) (*NsbGroupStateMachineService, error) {
	scenario, err := service.ChatScenarioService.GetScenario("nsbGroup")
	if err != nil {
		return nil, fmt.Errorf("error getting scenario: %w", err)
	}

	eviRole, err := service.ChatScenarioService.GetScenarioRole("nsbGroup", "evi")
	if err != nil {
		return nil, fmt.Errorf("error getting evi role: %w", err)
	}
	nsbGroupEviRoleStateMachine := NewNsbGroupEviRoleStateMachineService(eviRole, service)

	captainRole, err := service.ChatScenarioService.GetScenarioRole("nsbGroup", "captain")
	if err != nil {
		return nil, fmt.Errorf("error getting evi role: %w", err)
	}
	nsbGroupCaptainRoleStateMachine := NewNsbGroupCaptainRoleStateMachineService(captainRole, service)

	reviewerRole, err := service.ChatScenarioService.GetScenarioRole("nsbGroup", "reviewer")
	if err != nil {
		return nil, fmt.Errorf("error getting evi role: %w", err)
	}
	nsbGroupReviewerRoleStateMachine := NewNsbGroupReviewerRoleStateMachineService(reviewerRole, service)

	return &NsbGroupStateMachineService{
		Scenario:                        scenario,
		StateMachineService:             service,
		EviRoleStateMachineService:      nsbGroupEviRoleStateMachine,
		CaptainRoleStateMachineService:  nsbGroupCaptainRoleStateMachine,
		ReviewerRoleStateMachineService: nsbGroupReviewerRoleStateMachine}, nil
}

func (service *NsbGroupStateMachineService) GetName() string {
	return service.Scenario.Name
}

func (service *NsbGroupStateMachineService) ExecuteScenarioStateTransition(chatId string, req *models.ExecuteTransitionReq) (*models.State, error) {
	var nextState *models.State
	var err error
	switch req.RoleAction.ActiveRoleName {
	case "evi":
		nextState, _ = service.EviRoleStateMachineService.ExecuteScenarioRoleStateTransition(chatId, req)
	case "captain":
		nextState, err = service.CaptainRoleStateMachineService.ExecuteScenarioRoleStateTransition(chatId, req)
	case "reviewer":
		nextState, _ = service.ReviewerRoleStateMachineService.ExecuteScenarioRoleStateTransition(chatId, req)
	default:
		return nil, fmt.Errorf("no '%s' role found for scenario '%s'", req.RoleAction.ActiveRoleName, service.Scenario.Name)
	}
	if err != nil {
		return nil, err
	}
	return nextState, nil
}
