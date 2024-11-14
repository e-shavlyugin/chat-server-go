package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		HTTP
		Log
		Auth
		ChatContext
		ChatScenarioState
	}

	App struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Environment string `env:"ENVIRONMENT_NAME"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT" env-default:":3001"`

		Cors []string `yaml:"cors"`
	}

	Log struct {
		Level string `yaml:"level"`
	}

	Auth struct {
		AccessToken    string `env:"ACCESS_TOKEN"`
		SendgridToken  string `env:"SENDGRID_API_TOKEN"`
		GoogleAccessID string `env:"GOOGLE_ACCESS_ID"`
	}

	Cloud struct {
		BucketName string `env:"CLOUD_BUCKED_NAME" env-default:"evist-videos"`
	}

	Integrations struct {
		ChatServerEndpoint string `yaml:"chat-server-endpoint"`
	}

	ChatContext struct {
		ChatContextRepoType string `env:"CHAT_CONTEXT_REPO_TYPE" env-default:"mongodb"`
		MongoDB             struct {
			URI        string `env:"MONGODB_URI" env-default:"mongodb://mongodb:27017"`
			DB         string `env:"MONGODB_DATABASE"  env-default:"EVIDataStorage"`
			Collection string `env:"MONGODB_COLLECTION"  env-default:"ChatContext"`
		} `yaml:"mongodb"`
		File string `env:"PATH" env-default:"???"`
	}

	ChatScenarioState struct {
		ChatContextRepoType string `env:"CHAT_ROLE_STATE_REPO_TYPE" env-default:"mongodb"`
		MongoDB             struct {
			URI        string `env:"MONGODB_URI" env-default:"mongodb://mongodb:27017"`
			DB         string `env:"MONGODB_DATABASE"  env-default:"EVIDataStorage"`
			Collection string `env:"MONGODB_COLLECTION"  env-default:"ChatScenarioState"`
		}
		File string `env:"PATH" env-default:"???"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
