package sources

import (
	"chat_server_v2/internal/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ChatScenarioFileSource struct {
	Scenarios []models.ScenarioObject
}

func NewChatScenarioFileSource() (*ChatScenarioFileSource, error) {
	// Read the file
	var scenariosFilePath = "internal/repositories/sources/chat_scenarios.yaml"
	yamlFile, err := os.ReadFile(scenariosFilePath) // Fixed argument to use variable
	if err != nil {
		return &ChatScenarioFileSource{}, fmt.Errorf("error reading YAML file: %w", err)
	}

	// Declare an instance of ScenarioObject to hold the deserialized data
	var scenarios []models.ScenarioObject

	// Unmarshal YAML into the ScenarioObject struct
	err = yaml.Unmarshal(yamlFile, &scenarios)
	if err != nil {
		return &ChatScenarioFileSource{}, fmt.Errorf("error deserializing YAML: %w", err)
	}
	fmt.Println("yoyoyo1")
	// Return the deserialized struct and nil error if successful
	return &ChatScenarioFileSource{Scenarios: scenarios}, nil
}

func (source *ChatScenarioFileSource) GetScenario(name string) (*models.ScenarioObject, error) {
	// Loop through all roles and find the one with the provided name
	fmt.Println("yoyoyo")
	for _, scenario := range source.Scenarios {
		if scenario.Name == name {
			return &scenario, nil
		}
	}
	// If no matching role is found, return an error
	return nil, fmt.Errorf("scenario with name '%s' not found", name)
}

func (source *ChatScenarioFileSource) ListScenarios() ([]models.ScenarioObject, error) {
	return source.Scenarios, nil
}
