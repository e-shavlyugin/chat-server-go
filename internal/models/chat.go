package models

import (
	"github.com/tmc/langchaingo/llms"
	"time"
)

type ChatContextObj struct {
	ID           string  `json:"id" bson:"_id,omitempty"`
	ScenarioName string  `json:"scenarioName"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Events       []Event `json:"events"`
	Count        int     `json:"count"`
}

type Event struct {
	ID      int          `json:"id"`
	Role    string       `json:"role"` // User or AI
	Type    EventType    `json:"type"` // prompt or completion
	Payload EventPayload `json:"payload"`
	Time    time.Time    `json:"time"`
}

type EventType string

const (
	EventTypeSendPrompt          EventType = "sendPrompt"
	EventTypeGetCompletion       EventType = "getCompletion"
	EventTypeApplySystemTemplate EventType = "applySystemTemplate"
	EventTypeTodo                EventType = "eventTypeTodo"
)

type EventPayload struct {
	AnyOtherTypeOfMessage string
	ChatMessage           ChatMessage
}

type ChatMessage struct {
	Message         string
	MessageExpected string               // for llm tuning
	LLMRole         llms.ChatMessageType `json:"llmRole"`
	Rating          LikertRatingType     `json:"rating"` // Likert scale rating (1-5)
}

type LikertRatingType int

const (
	VeryUnsatisfied                 LikertRatingType = 1
	UnsatisfiedLikertRating         LikertRatingType = 2
	NeitherSatisfiedNorDissatisfied LikertRatingType = 3
	Satisfied                       LikertRatingType = 4
	VerySatisfied                   LikertRatingType = 5
)

type CreateChatContextObjectReq struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ScenarioName string `json:"scenarioName"`
}
