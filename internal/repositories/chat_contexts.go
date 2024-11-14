package repositories

import (
	"chat_server_v2/config"
	"chat_server_v2/internal/models"
	"chat_server_v2/internal/repositories/sources"
	"log"
)

type ChatContextsSource interface {
	CreateChatContext(chatContextObject *models.ChatContextObj) (*models.ChatContextObj, error)
	GetChatContext(id string) (*models.ChatContextObj, error)
	ListChatContexts() ([]models.ChatContextObj, error)
	UpdateChatContextEvents(id string, events []models.Event) (string, error)
	SaveChat(chat *models.ChatContextObj) (string, error)
	DeleteChat(id string) (string, error)
}

type ChatContextRepository struct {
	chatContextSource ChatContextsSource
}

func NewChatContextsRepository(cfg *config.Config) *ChatContextRepository {

	var chatContextSource ChatContextsSource

	switch cfg.ChatContext.ChatContextRepoType {
	case "http":
		// chatRepo = repository.NewChatRepository("http://example.com") // Replace with your service URL
	case "mongodb":
		chatContextSource = sources.NewMongoDBChatContextSource(cfg.ChatContext.MongoDB.URI, cfg.ChatContext.MongoDB.DB, cfg.ChatContext.MongoDB.Collection)
	default:
		log.Fatalf("Unsupported repository type: %s", cfg.ChatContext.ChatContextRepoType)
	}

	return &ChatContextRepository{chatContextSource: chatContextSource}

}

func (r *ChatContextRepository) CreateChatContext(chatContextObject *models.ChatContextObj) (*models.ChatContextObj, error) {
	return r.chatContextSource.CreateChatContext(chatContextObject)
}

func (r *ChatContextRepository) GetChatContext(chatId string) (*models.ChatContextObj, error) {
	return r.chatContextSource.GetChatContext(chatId)
}

func (r *ChatContextRepository) ListChatContexts() ([]models.ChatContextObj, error) {
	return r.chatContextSource.ListChatContexts()
}

func (r *ChatContextRepository) UpdateChatContextEvents(chatId string, events []models.Event) (string, error) {
	return r.chatContextSource.UpdateChatContextEvents(chatId, events)
}

func (r *ChatContextRepository) SaveChat(chat *models.ChatContextObj) (string, error) {
	return r.chatContextSource.SaveChat(chat)
}
