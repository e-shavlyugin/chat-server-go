package chat

import (
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/repositories"
	//"chat_server_v2/internal/services"
)

type ChatScenarioStateService struct {
	ChatScenarioStatesRepository *repositories.ChatScenarioStatesRepository
	ChatContextService           *ChatContextService
}

func NewChatScenarioStateService(repository *repositories.ChatScenarioStatesRepository, service *ChatContextService) *ChatScenarioStateService {
	return &ChatScenarioStateService{ChatScenarioStatesRepository: repository, ChatContextService: service}
}

func (service *ChatScenarioStateService) CreateChatScenarioState(obj *models.ChatScenarioStateObj) (*models.ChatScenarioStateObj, error) {
	return service.ChatScenarioStatesRepository.CreateChatScenarioState(obj)
}

func (service *ChatScenarioStateService) GetChatScenarioState(chatId string) (*models.ChatScenarioStateObj, error) {
	return service.ChatScenarioStatesRepository.GetChatScenarioState(chatId)
}

func (service *ChatScenarioStateService) ListChatScenarioStates() ([]models.ChatScenarioStateObj, error) {
	return service.ChatScenarioStatesRepository.ListChatScenarioStates()
}

func (service *ChatScenarioStateService) UpdateChatScenarioState(obj *models.ChatScenarioStateObj) (*models.ChatScenarioStateObj, error) {
	return service.ChatScenarioStatesRepository.UpdateChatScenarioState(obj)
}
