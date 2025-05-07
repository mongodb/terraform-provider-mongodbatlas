package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
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

const DefaultAPIVersion = "2024-12-01-preview"
const Model = "gpt-4.1"

func CallModel(client *openai.Client, systemPrompt string, userPrompt string) string {
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model:       Model,
		Temperature: param.NewOpt(0.0),
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content
}

// GetRequiredEnvVar retrieves an environment variable and returns an error if it's not set
func GetRequiredEnvVar(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", name)
	}
	return value, nil
}
