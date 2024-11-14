package models

type InitStateMachineReq struct {
	ScenarioName string `json:"scenarioName"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}
