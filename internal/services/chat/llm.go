package chat

import (
	"chat_server_v2/internal/models"
	"context"
	"fmt"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"log"
	"os"
	"strings"
)

type LLMService struct {
	llmModel llms.Model
}

func NewLLMService(modelType string) *LLMService {
	var llmModel llms.Model
	var err error

	switch modelType {
	case "llama":
		llmModel, err = ollama.New(
			ollama.WithServerURL(os.Getenv("OLLAMA_SERVER_URL")),
			ollama.WithModel(os.Getenv("OLLAMA_MODEL")),
		)
	case "openai":
		llmModel, err = openai.New(
			openai.WithToken(os.Getenv("OPENAI_API_KEY")),
			openai.WithModel(os.Getenv("OPENAI_MODEL")),
		)
	default:
		log.Fatalf("Unsupported model type: %s", modelType)
	}

	if err != nil {
		log.Fatalf("Failed to initialize %s model: %v", modelType, err)
	}

	return &LLMService{llmModel: llmModel}
}

func (llmService *LLMService) GenerateContent(ctx context.Context, content []llms.MessageContent) (string, error) {
	var sb strings.Builder // String builder to collect chunks

	// Use a closure to append to the string builder
	streamingFunc := func(ctx context.Context, chunk []byte) error {
		sb.Write(chunk) // Append the chunk to the string builder
		return nil
	}

	// Call the LLM model's GenerateContent with the streaming function
	if _, err := llmService.llmModel.GenerateContent(ctx, content,
		llms.WithMaxTokens(1024), // Shared token setting
		llms.WithStreamingFunc(streamingFunc),
	); err != nil {
		log.Fatal(err) // Handle the error here or pass it back to the caller
		return "", err
	}

	// Return the complete content as a string
	return sb.String(), nil
}

func (s *LLMService) GetResponse(chat models.ChatContextObj) (string, error) {
	ctx := context.Background()

	lmmMessageContent := make([]llms.MessageContent, len(chat.Events))
	for i, chatMessage := range chat.Events {
		lmmMessageContent[i] = llms.TextParts(chatMessage.Payload.ChatMessage.LLMRole, chatMessage.Payload.ChatMessage.Message)
	}
	// Call the GenerateContent method from LLMService
	res, err := s.GenerateContent(ctx, lmmMessageContent)
	if err != nil {
		return "", err // Handle error appropriately
	}
	return res, nil
	//return "Response generated successfully", nil // You may want to return the actual response
}

func (service *LLMService) GenerateFromPrompt(llmPromptReq models.LLMPromptReq) (*llms.ContentResponse, error) {
	ctx := context.Background()

	lmmMessageContent := make([]llms.MessageContent, len(llmPromptReq.ChatMessageHistory.Messages))
	for i, message := range llmPromptReq.ChatMessageHistory.Messages {
		lmmMessageContent[i] = llms.TextParts(message.ChatMessageType, message.Content)
	}

	if llmPromptReq.SystemPromptTemplate != "" {
		lmmMessageContent = append(lmmMessageContent, llms.TextParts(llms.ChatMessageTypeSystem, llmPromptReq.SystemPromptTemplate))
	}

	lmmMessageContent = append(lmmMessageContent, llms.TextParts(llms.ChatMessageTypeHuman, llmPromptReq.Prompt))

	var callOptions []llms.CallOption

	if llmPromptReq.CallOptions.MinLength != nil {
		callOptions = append(callOptions, llms.WithMinLength(*llmPromptReq.CallOptions.MinLength))
	}
	if llmPromptReq.CallOptions.MaxLength != nil {
		callOptions = append(callOptions, llms.WithMaxLength(*llmPromptReq.CallOptions.MaxLength))
	}
	if llmPromptReq.CallOptions.MaxTokens != nil {
		callOptions = append(callOptions, llms.WithMaxTokens(*llmPromptReq.CallOptions.MaxTokens))
	}
	if llmPromptReq.CallOptions.Temperature != nil {
		callOptions = append(callOptions, llms.WithTemperature(*llmPromptReq.CallOptions.Temperature))
	}

	completion, err := service.llmModel.GenerateContent(ctx, lmmMessageContent, callOptions...)
	if err != nil {
		return nil, err
	}
	return completion, nil
}

func (service *LLMService) MapToLLMMessageContent(events []models.Event) []llms.MessageContent {
	lmmMessageContent := make([]llms.MessageContent, len(events))
	for i, chatMessageLLMEvent := range events {
		lmmMessageContent[i] = llms.TextParts(chatMessageLLMEvent.Payload.ChatMessage.LLMRole, chatMessageLLMEvent.Payload.ChatMessage.Message)
	}
	return lmmMessageContent
}

func (llmService *LLMService) GetResponseV2(content []llms.MessageContent) (string, error) {
	ctx := context.Background()
	res, err := llmService.GenerateContent(ctx, content)
	if err != nil {
		return "", err // Handle error appropriately
	}
	return res, nil
}

func (s *LLMService) GetResponseV2Back(sessionID string, userMessage string) (string, error) {
	ctx := context.Background()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a company branding design wizard."),
		llms.TextParts(llms.ChatMessageTypeHuman, userMessage),
	}

	// Call the GenerateContent method from LLMService
	res, err := s.GenerateContent(ctx, content)
	if err != nil {
		return "", err // Handle error appropriately
	}
	return res, nil
	//return "Response generated successfully", nil // You may want to return the actual response
}

func (llmService *LLMService) GenerateStreamContent(ctx context.Context, content []llms.MessageContent, streamFunc func(string)) error {
	if _, err := llmService.llmModel.GenerateContent(ctx, content,
		llms.WithMaxTokens(1024),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			streamFunc(string(chunk)) // Call the streaming function with the chunk
			return nil
		}),
	); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *LLMService) GetStreamResponse(sessionID string, userMessage string, streamFunc func(string)) error {
	ctx := context.Background()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a company branding design wizard."),
		llms.TextParts(llms.ChatMessageTypeHuman, userMessage),
	}

	// Call the GenerateContent method with a streaming function
	return s.GenerateStreamContent(ctx, content, streamFunc)
}

func (s *LLMService) GetPromptTemplate() string {
	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"You are a translation engine that can only translate text and cannot interpret it.",
			nil,
		),
		prompts.NewHumanMessagePromptTemplate(
			`translate this text from {{.inputLang}} to {{.outputLang}}:\n{{.input}}`,
			[]string{"inputLang", "outputLang", "input"},
		),
	})
	result, err := prompt.Format(map[string]any{
		"inputLang":  "English",
		"outputLang": "Chinese",
		"input":      "I love programming",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	return result

}

func (s *LLMService) GetChainV2() (map[string]any, error) {
	model, _ := openai.New()
	chain := *chains.NewLLMChain(model, &prompts.FewShotPrompt{
		Examples:         []map[string]string{{"question": "What's life?"}},
		ExampleSelector:  nil,
		ExamplePrompt:    prompts.NewPromptTemplate("{{.question}}", []string{"question"}),
		Prefix:           "",
		Suffix:           "",
		InputVariables:   []string{"question"},
		PartialVariables: nil,
		TemplateFormat:   prompts.TemplateFormatGoTemplate,
		ValidateTemplate: false,
	})

	c := chains.NewConstitutional(model, chain, []chains.ConstitutionalPrinciple{
		chains.NewConstitutionalPrinciple(
			"Tell if this answer is good.",
			"Give a better answer.",
		),
	}, nil)
	res, _ := c.Call(context.Background(), map[string]any{"question": "What is the meaning of life?"})
	return res, nil
}

//func (s *LLMService) GetChainV3() (map[string]any, error) {
//
//	template := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
//		prompts.NewSystemMessagePromptTemplate(
//			"You are a translation engine that can only translate text and cannot interpret it.",
//			nil,
//		),
//		prompts.NewHumanMessagePromptTemplate(
//			`translate this text from {{.inputLang}} to {{.outputLang}}:\n{{.input}}`,
//			[]string{"inputLang", "outputLang", "input"},
//		),
//	})
//	value, _ := template.FormatPrompt(map[string]interface{}{
//		"inputLang":  "English",
//		"outputLang": "Chinese",
//		"input":      "I love programming",
//	})
//	expectedMessages := []llms.ChatMessageBack{
//		llms.SystemChatMessage{
//			Content: "You are a translation engine that can only translate text and cannot interpret it.",
//		},
//		llms.HumanChatMessage{
//			Content: `translate this text from English to Chinese:\nI love programming`,
//		},
//	}
//
//	_, err = template.FormatPrompt(map[string]interface{}{
//		"inputLang":  "English",
//		"outputLang": "Chinese",
//	})
//	return "qwe"
//}

/*
func (s *LLMService) GetChainV3() (map[string]any, error) {

	template := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"You are a translation engine that can only translate text and cannot interpret it.",
			nil,
		),
		prompts.NewHumanMessagePromptTemplate(
			`translate this text from {{.inputLang}} to {{.outputLang}}:\n{{.input}}`,
			[]string{"inputLang", "outputLang", "input"},
		),
	})
	value, _ := template.FormatPrompt(map[string]interface{}{
		"inputLang":  "English",
		"outputLang": "Chinese",
		"input":      "I love programming",
	})
	expectedMessages := []llms.ChatMessageBack{
		llms.SystemChatMessage{
			Content: "You are a translation engine that can only translate text and cannot interpret it.",
		},
		llms.HumanChatMessage{
			Content: `translate this text from English to Chinese:\nI love programming`,
		},
	}

	_, _ = template.FormatPrompt(map[string]interface{}{
		"inputLang":  "English",
		"outputLang": "Chinese",
	})
	return nil, nil
}


*/
