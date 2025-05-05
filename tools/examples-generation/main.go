package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	// "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	// "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
)

const atlasAdminAPISpecURL = "https://raw.githubusercontent.com/mongodb/openapi/refs/heads/main/openapi/v2/openapi-2025-03-12.yaml"
const specFilePath = "open-api-spec.yml"

func main() {

	// if err := openapi.DownloadOpenAPISpec(atlasAdminAPISpecURL, specFilePath); err != nil {
	// 	log.Fatalf("an error occurred when downloading Atlas Admin API spec: %v", err)
	// }

	// apiSpec, err := openapi.ParseAtlasAdminAPI(specFilePath)
	// if err != nil {
	// 	log.Fatalf("unable to parse Atlas Admin API: %v", err)
	// }

	// op, _ := apiSpec.Model.Paths.PathItems.OrderedMap.Get("/api/atlas/v2/groups/{groupId}/alertConfigs/{alertConfigId}")

	// okResponse, _ := op.Get.Responses.Codes.Get(codespec.OASResponseCodeOK)
	// schema, _ := codespec.GetSchemaFromMediaType(okResponse.Content)
	// baseSchema := schema.Schema
	// yamlBytes, _ := baseSchema.RenderInline()
	// yamlString := string(yamlBytes[:])
	// println(yamlString)
	// print(schema)

	const DefaultAPIVersion = "2024-12-01-preview"
	const Model = "gpt-4.1-nano"

	client, err := CreateOpenAIClientWithKey(DefaultAPIVersion)
	if err != nil {
		panic(fmt.Errorf("failed to create OpenAI client: %w", err))
	}

	// Read the system prompt
	systemPromptBytes, err := os.ReadFile("tools/examples-generation/prompts/generatehcl.system.md")
	if err != nil {
		panic(fmt.Errorf("failed to read system prompt file: %w", err))
	}
	systemPrompt := string(systemPromptBytes)

	// Read the user prompt
	userPromptBytes, err := os.ReadFile("tools/examples-generation/prompts/generatehcl.user.md")
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
