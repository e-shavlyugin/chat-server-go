package models

import "github.com/tmc/langchaingo/llms"

type ScenarioObject struct {
	Name  string `yaml:"name"`
	Roles []Role `yaml:"roles"`
}

type Role struct {
	Name               string               `yaml:"name"`
	Description        string               `yaml:"description,omitempty"`
	LLMChatMessageType llms.ChatMessageType `yaml:"LLMChatMessageType"`
	Actions            []Action             `yaml:"actions"`
}

type Action struct {
	Name                   string    `yaml:"name"`
	Description            string    `yaml:"description"`
	NextRole               string    `yaml:"nextRole"`
	NextGuaranteedActions  []string  `yaml:"nextGuaranteedActions"`
	NextConditionalActions []string  `yaml:"nextConditionalActions"`
	LLMConfig              LLMConfig `yaml:"LLMConfig"`
}

type LLMConfig struct {
	MaxTokens            int    `yaml:"maxTokens"`
	SystemPromptTemplate string `yaml:"systemPromptTemplate"`
}
