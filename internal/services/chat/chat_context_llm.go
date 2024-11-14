package chat

import (
	"chat_server_v2/internal/models"
	"github.com/tmc/langchaingo/llms"
	"time"
)

type ChatContextLLMService struct {
	ChatContextService *ChatContextService
	LLMService         *LLMService
}

func NewChatContextLLMService(chatContextService *ChatContextService, llmService *LLMService) *ChatContextLLMService {
	return &ChatContextLLMService{LLMService: llmService, ChatContextService: chatContextService}
}

func (chatService *ChatContextLLMService) GetLLMChatContextHistory(chatId string) ([]models.Event, error) {
	chatContextObject, err := chatService.ChatContextService.GetChatContext(chatId)
	if err != nil {
		return nil, err
	}
	return chatService.retrieveLLMChatContextHistory(chatContextObject), nil
}

func (chatService *ChatContextLLMService) retrieveLLMChatContextHistory(chatContextObject *models.ChatContextObj) []models.Event {
	var chatMessageLLMEvents []models.Event
	var chatMessageEventsPredicate = func(e models.Event) bool {
		return e.Type == models.EventTypeSendPrompt ||
			e.Type == models.EventTypeGetCompletion ||
			e.Type == models.EventTypeApplySystemTemplate
	}
	for _, chatMessageLLMEvent := range chatContextObject.Events {
		if chatMessageEventsPredicate(chatMessageLLMEvent) {
			chatMessageLLMEvents = append(chatMessageLLMEvents, chatMessageLLMEvent)
		}
	}
	return chatMessageLLMEvents
}

func (chatService *ChatContextLLMService) MapToLLMMessageContent(events []models.Event) []llms.MessageContent {
	lmmMessageContent := make([]llms.MessageContent, len(events))
	for i, chatMessageLLMEvent := range events {
		lmmMessageContent[i] = llms.TextParts(chatMessageLLMEvent.Payload.ChatMessage.LLMRole, chatMessageLLMEvent.Payload.ChatMessage.Message)
	}
	return lmmMessageContent
}

// / ???
func (chatService *ChatContextLLMService) SendPromptTemp(chatId string, roleName string, systemPromptTemplate string, prompt string) (*models.ChatContextObj, error) {
	// 1. get chat
	chatContextObject, err := chatService.ChatContextService.GetChatContext(chatId)
	if err != nil {
		return nil, err
	}
	return chatContextObject, nil
}

func (service *ChatContextLLMService) SendPrompt(chatContextObject *models.ChatContextObj, roleName string, systemPromptTemplate string, prompt string) (*models.ChatContextObj, error) {

	var counter = chatContextObject.Count
	if systemPromptTemplate != "" && service.GetLatestSystemPromptTemplate(chatContextObject) != systemPromptTemplate {
		counter++
		newSystemPromptTemplateMessage := models.Event{
			ID:   counter,
			Role: roleName,
			Type: models.EventTypeApplySystemTemplate,
			Payload: models.EventPayload{
				ChatMessage: models.ChatMessage{
					Message: systemPromptTemplate,
					LLMRole: llms.ChatMessageTypeSystem,
				},
			},
			Time: time.Now(),
		}
		chatContextObject.Events = append(chatContextObject.Events, newSystemPromptTemplateMessage)
	}

	if prompt != "" {
		counter++
		newHumanMessage := models.Event{
			ID:   counter,
			Role: roleName,
			Type: models.EventTypeSendPrompt,
			Payload: models.EventPayload{
				ChatMessage: models.ChatMessage{
					Message: prompt,
					LLMRole: llms.ChatMessageTypeHuman,
				},
			},
			Time: time.Now(),
		}
		chatContextObject.Events = append(chatContextObject.Events, newHumanMessage)

		llmResponse, err := service.LLMService.GetResponseV2(
			service.MapToLLMMessageContent(
				service.retrieveLLMChatContextHistory(chatContextObject)))
		if err != nil {
			return nil, err
		}
		newAIMessage := models.Event{
			ID:   chatContextObject.Count + 3,
			Role: roleName,
			Type: models.EventTypeGetCompletion,
			Payload: models.EventPayload{
				ChatMessage: models.ChatMessage{
					Message: llmResponse,
					LLMRole: llms.ChatMessageTypeAI,
				},
			},
			Time: time.Now(),
		}
		chatContextObject.Events = append(chatContextObject.Events, newAIMessage)
	}
	_, err := service.ChatContextService.UpdateChatContextEvents(chatContextObject.ID, chatContextObject.Events)
	if err != nil {
		return nil, err
	}
	return chatContextObject, nil
}

func (service *ChatContextLLMService) SendAIMessage(chatContextObject *models.ChatContextObj, roleName string, systemPromptTemplate string, message string) (*models.ChatContextObj, error) {

	var counter = chatContextObject.Count
	if systemPromptTemplate != "" && service.GetLatestSystemPromptTemplate(chatContextObject) != systemPromptTemplate {
		counter++
		newSystemPromptTemplateMessage := models.Event{
			ID:   counter,
			Role: roleName,
			Type: models.EventTypeApplySystemTemplate,
			Payload: models.EventPayload{
				ChatMessage: models.ChatMessage{
					Message: systemPromptTemplate,
					LLMRole: llms.ChatMessageTypeSystem,
				},
			},
			Time: time.Now(),
		}
		chatContextObject.Events = append(chatContextObject.Events, newSystemPromptTemplateMessage)
	}

	if message != "" {
		counter++
		newAIMessage := models.Event{
			ID:   counter,
			Role: roleName,
			Type: models.EventTypeSendPrompt,
			Payload: models.EventPayload{
				ChatMessage: models.ChatMessage{
					Message: message,
					LLMRole: llms.ChatMessageTypeAI,
				},
			},
			Time: time.Now(),
		}
		chatContextObject.Events = append(chatContextObject.Events, newAIMessage)

		response, err := service.LLMService.GetResponseV2(
			service.MapToLLMMessageContent(
				service.retrieveLLMChatContextHistory(chatContextObject)))
		println(response)
		if err != nil {
			return nil, err
		}
	}
	_, err := service.ChatContextService.UpdateChatContextEvents(chatContextObject.ID, chatContextObject.Events)
	if err != nil {
		return nil, err
	}
	return chatContextObject, nil
}

func (service *ChatContextLLMService) GetLatestSystemPromptTemplate(chatContextObject *models.ChatContextObj) string {
	var latest *models.Event

	// Iterate over the array and find the latest green apple
	for _, event := range chatContextObject.Events {
		if event.Type == models.EventTypeApplySystemTemplate {
			latest = &event
		}
	}

	if latest == nil {
		return ""
	}
	return latest.Payload.ChatMessage.Message
}

/*
func (service *ChatContextLLMService) SendSystemMessage(chatContextObject *models.ChatContextObj, roleName string, systemPromptTemplate string, prompt string) (*models.ChatContextObj, error) {

	var counter = chatContextObject.Count
	if systemPromptTemplate != "" {
		counter++
		newSystemPromptTemplateMessage := models.Event{
			ID:        counter,
			Role:      roleName,
			Type: models.EventTypeApplySystemTemplate,
			Payload: models.EventPayload{
				ChatMessage: models.ChatMessage{
					Message: systemPromptTemplate,
					LLMRole: llms.ChatMessageTypeSystem,
				},
			},
			Time: time.Now(),
		}
		chatContextObject.Events = append(chatContextObject.Events, newSystemPromptTemplateMessage)
	}

	if prompt != "" {
		counter++
		newHumanMessage := models.Event{
			ID:        counter,
			Role:      roleName,
			Type: models.EventTypeSendPrompt,
			Payload: models.EventPayload{
				ChatMessage: models.ChatMessage{
					Message: prompt,
					LLMRole: llms.ChatMessageTypeHuman,
				},
			},
			Time: time.Now(),
		}
		chatContextObject.Events = append(chatContextObject.Events, newHumanMessage)

		llmResponse, err := service.LLMService.GetResponseV2(
			service.MapToLLMMessageContent(
				service.retrieveLLMChatContextHistory(chatContextObject)))
		if err != nil {
			return nil, err
		}
		newAIMessage := models.Event{
			ID:        chatContextObject.Count + 3,
			Role:      roleName,
			Type: models.EventTypeGetCompletion,
			Payload: models.EventPayload{
				ChatMessage: models.ChatMessage{
					Message: llmResponse,
					LLMRole: llms.ChatMessageTypeAI,
				},
			},
			Time: time.Now(),
		}
		chatContextObject.Events = append(chatContextObject.Events, newAIMessage)
	}
	_, err := service.ChatContextService.UpdateChatContextEvents(chatContextObject.ID, chatContextObject.Events)
	if err != nil {
		return nil, err
	}
	return chatContextObject, nil
}
*/
