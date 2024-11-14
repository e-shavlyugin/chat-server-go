package repositories

import (
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/repositories/sources"
)

type ChatScenarioSource interface {
	GetScenario(name string) (*models.ScenarioObject, error)
	ListScenarios() ([]models.ScenarioObject, error)
}

type ChatScenarioRepository struct {
	Source ChatScenarioSource
}

func NewChatScenarioRepository(scenarioRoleSourceType string) (*ChatScenarioRepository, error) {
	var source ChatScenarioSource
	var err error
	switch scenarioRoleSourceType {
	case "file":
		source, err = sources.NewChatScenarioFileSource()
	case "mongodb":
		source = nil
	}
	if err != nil {
		return &ChatScenarioRepository{}, err
	}
	return &ChatScenarioRepository{Source: source}, nil
}

func (r *ChatScenarioRepository) GetScenario(name string) (*models.ScenarioObject, error) {
	return r.Source.GetScenario(name)
}

func (r *ChatScenarioRepository) ListScenarios() ([]models.ScenarioObject, error) {
	return r.Source.ListScenarios()
}
