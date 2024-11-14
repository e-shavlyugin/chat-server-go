package repositories

import (
	"chat_server_v2/config"
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/repositories/sources"
)

type ChatScenarioStateSource interface {
	CreateChatScenarioState(obj *models.ChatScenarioStateObj) (*models.ChatScenarioStateObj, error)
	GetChatScenarioState(chatId string) (*models.ChatScenarioStateObj, error)
	ListChatScenarioStates() ([]models.ChatScenarioStateObj, error)
	UpdateChatScenarioState(obj *models.ChatScenarioStateObj) (*models.ChatScenarioStateObj, error)
}

type ChatScenarioStatesRepository struct {
	DialogueStateSource ChatScenarioStateSource
}

func NewChatScenarioStateRepository(cfg *config.Config) (*ChatScenarioStatesRepository, error) {
	var dialogueStateSource ChatScenarioStateSource
	var err error
	switch cfg.ChatScenarioState.ChatContextRepoType {
	case "mongodb":
		dialogueStateSource, err = sources.NewMongoDBChatStateSource(cfg.ChatScenarioState.MongoDB.URI, cfg.ChatScenarioState.MongoDB.DB, cfg.ChatScenarioState.MongoDB.Collection)
	case "other":
		dialogueStateSource, err = nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ChatScenarioStatesRepository{DialogueStateSource: dialogueStateSource}, nil
}

func (r *ChatScenarioStatesRepository) CreateChatScenarioState(obj *models.ChatScenarioStateObj) (*models.ChatScenarioStateObj, error) {
	return r.DialogueStateSource.CreateChatScenarioState(obj)
}

func (r *ChatScenarioStatesRepository) GetChatScenarioState(chatId string) (*models.ChatScenarioStateObj, error) {
	return r.DialogueStateSource.GetChatScenarioState(chatId)
}

func (r *ChatScenarioStatesRepository) ListChatScenarioStates() ([]models.ChatScenarioStateObj, error) {
	return r.DialogueStateSource.ListChatScenarioStates()
}

func (r *ChatScenarioStatesRepository) UpdateChatScenarioState(dialogueStateObject *models.ChatScenarioStateObj) (*models.ChatScenarioStateObj, error) {
	return r.DialogueStateSource.UpdateChatScenarioState(dialogueStateObject)
}
