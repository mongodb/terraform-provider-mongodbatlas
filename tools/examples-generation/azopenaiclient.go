package main

import (
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"github.com/openai/openai-go/option"
)

func CreateOpenAIClientWithKey(apiVersion string) (*openai.Client, error) {
	apiKey, err := GetRequiredEnvVar("AZURE_OPENAI_API_KEY")
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	endpoint, err := GetRequiredEnvVar("AZURE_OPENAI_ENDPOINT")
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint: %w", err)
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		azure.WithEndpoint(endpoint, apiVersion),
	)
	return &client, nil
}

// GetRequiredEnvVar retrieves an environment variable and returns an error if it's not set
func GetRequiredEnvVar(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", name)
	}
	return value, nil
}
