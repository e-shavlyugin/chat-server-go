package models

type ChatScenarioStateObj struct {
	ChatId string `json:"chatId" bson:"_id,omitempty"`
	State  State  `json:"state"`
}

type State struct {
	NextActiveRoleName string
	NextActions        []string
}

type RoleAction struct {
	ActiveRoleName string `json:"activeRoleName"`
	ActionName     string `json:"actionName"`
}

/*
	type ExecuteTransitionReq struct {
		ScenarioName string                   `json:"scenarioName"`
		ChatId       string                   `json:"chatId"`
		Payload      TransitionRequestPayload `json:"payload"`
		RoleAction   RoleAction               `json:"roleAction"`
	}
*/
type ExecuteTransitionReq struct {
	Payload    TransitionRequestPayload `json:"payload"`
	RoleAction RoleAction               `json:"roleAction"`
}

type TransitionRequestPayload struct {
	SomethingElse string `json:"somethingElse"`
	Message       string `json:"message"`
}

type CreateChatScenarioStateReq struct {
	ScenarioName string `json:"scenarioName"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}
