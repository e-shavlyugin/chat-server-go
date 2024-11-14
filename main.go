package main

import (
	"chat_server_v2/config"
	"chat_server_v2/internal/logger"
	"chat_server_v2/internal/server"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// @title Chat Server API
// @version 1.0
// @description This is the API documentation for the Chat Server
// @BasePath /v1
func main() {

	l := logger.NewLogger()

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	cfg, err := config.NewConfig(
		fmt.Sprintf(
			"config/app.%s.yaml",
			os.Getenv("ENVIRONMENT_NAME"),
		),
	)
	if err != nil {
		l.Fatal("error", err)
	}

	if err = server.Start(cfg); err != nil {
		l.Fatal("error", err)
	}

}
