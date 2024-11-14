package models

import (
	"github.com/tmc/langchaingo/llms"
)

type LLMPromptReq struct {
	Prompt               string             `json:"prompt"`
	SystemPromptTemplate string             `json:"systemPromptTemplate"`
	CallOptions          CallOptions        `json:"callOptions"`
	ChatMessageHistory   ChatMessageHistory `json:"history"`
}

type ChatMessageHistory struct {
	Messages []Message
}

type Message struct {
	ChatMessageType llms.ChatMessageType `json:"type"`
	Content         string               `json:"content"`
}

type CallOptions struct {
	Temperature *float64 `json:"temperature"`
	MaxTokens   *int     `json:"maxTokens"`
	MinLength   *int     `json:"minLength"`
	MaxLength   *int     `json:"maxLength"`
}
