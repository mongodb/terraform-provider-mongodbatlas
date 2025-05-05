package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
)

func main() {
	const DefaultAPIVersion = "2024-12-01-preview"
	const Model = "gpt-4.1-nano"

	client, err := CreateOpenAIClientWithKey(DefaultAPIVersion)
	if err != nil {
		panic(fmt.Errorf("failed to create OpenAI client: %w", err))
	}

	// Read the system prompt
	systemPromptBytes, err := os.ReadFile("tools/examples-generation/prompts/generateexample.system.md")
	if err != nil {
		panic(fmt.Errorf("failed to read system prompt file: %w", err))
	}
	systemPrompt := string(systemPromptBytes)

	// Read the user prompt
	userPromptBytes, err := os.ReadFile("tools/examples-generation/prompts/generateexample.user.md")
	if err != nil {
		panic(fmt.Errorf("failed to read user prompt file: %w", err))
	}
	userPrompt := string(userPromptBytes)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model: Model,
	})

	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
}
