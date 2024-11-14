package chat

import (
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/repositories"
	"github.com/google/uuid"
)

type ChatContextService struct {
	ChatContextRepository *repositories.ChatContextRepository
}

func NewChatContextService(chatContextRepository *repositories.ChatContextRepository) *ChatContextService {
	return &ChatContextService{ChatContextRepository: chatContextRepository}
}

func (chatContextService *ChatContextService) CreateChatContext(obj *models.ChatContextObj) (*models.ChatContextObj, error) {
	obj.ID = uuid.New().String()
	return chatContextService.ChatContextRepository.CreateChatContext(obj)
}

func (chatContextService *ChatContextService) GetChatContext(chatId string) (*models.ChatContextObj, error) {
	return chatContextService.ChatContextRepository.GetChatContext(chatId)
}

func (chatContextService *ChatContextService) ListChatContexts() ([]models.ChatContextObj, error) {
	return chatContextService.ChatContextRepository.ListChatContexts()
}
func (chatContextService *ChatContextService) UpdateChatContextEvents(chatId string, events []models.Event) (string, error) {
	return chatContextService.ChatContextRepository.UpdateChatContextEvents(chatId, events)
}
