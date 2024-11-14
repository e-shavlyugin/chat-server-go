package chat

import (
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/repositories"
	"fmt"
)

type ChatScenarioService struct {
	chatScenarioRepository *repositories.ChatScenarioRepository
}

func NewChatScenarioService(repository *repositories.ChatScenarioRepository) *ChatScenarioService {
	return &ChatScenarioService{chatScenarioRepository: repository}
}

func (r *ChatScenarioService) GetScenario(name string) (*models.ScenarioObject, error) {
	return r.chatScenarioRepository.GetScenario(name)
}

func (r *ChatScenarioService) ListScenarios() ([]models.ScenarioObject, error) {
	return r.chatScenarioRepository.ListScenarios()
}

func (r *ChatScenarioService) GetScenarioRole(scenarioName string, roleName string) (*models.Role, error) {
	scenario, err := r.GetScenario(scenarioName)
	if err != nil {
		return nil, err
	}
	for _, role := range scenario.Roles {
		if role.Name == roleName {
			return &role, nil
		}
	}
	return nil, fmt.Errorf("role '%s' not found in scenario '%s'", roleName, scenarioName)
}
