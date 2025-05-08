package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/examples-generation/prompts"
	"github.com/openai/openai-go"
)

const atlasAdminAPISpecURL = "https://raw.githubusercontent.com/mongodb/openapi/refs/heads/main/openapi/v2/openapi-2025-03-12.yaml"

const specFilePath = "open-api-spec.yml"
const resourcesBasePath = "./internal/service/"

var resourceToGetPath = map[string]string{
	"alert_configuration":   "/api/atlas/v2/groups/{groupId}/alertConfigs/{alertConfigId}",
	"search_index":          "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes/{indexId}",
	"search_deployment":     "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/deployment",
	"stream_instance":       "/api/atlas/v2/groups/{groupId}/streams/{tenantName}",
	"push_based_log_export": "/api/atlas/v2/groups/{groupId}/pushBasedLogExport",
	"cluster":               "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}",
}

func main() {
	osArg := getOsArg()
	if(osArg == nil) {
		log.Fatalf("resource name is required")
	}
	resourceName := *osArg

	resourceGetAPISpec := getAPISpecSchema(resourceName)
	resourceImpl := getResourceImplementationSchema(resourceName)

	client, err := CreateOpenAIClientWithKey(DefaultAPIVersion)
	if err != nil {
		panic(fmt.Errorf("failed to create OpenAI client: %w", err))
	}

	mainHCL := GenerateMainHCL(client, resourceName, resourceImpl, resourceGetAPISpec)

	if err := writeContentToExamplesFolder(mainHCL, "main.tf", resourceName); err != nil {
		log.Fatalf("Error writing main.tf: %v", err)
	}
	if err := writeContentToExamplesFolder(providerHCLContent, "provider.tf", resourceName); err != nil {
		log.Fatalf("Error writing provider hcl to file: %v", err)
	}
	if err := writeContentToExamplesFolder(versionsHCLContent, "versions.tf", resourceName); err != nil {
		log.Fatalf("Error writing versions hcl to file: %v", err)
	}

	variablesDefinitionHCL := GenerateVariableDefsHCL(client, mainHCL)

	if err := writeContentToExamplesFolder(variablesDefinitionHCL, "variables.tf", resourceName); err != nil {
		log.Fatalf("Error writing main.tf: %v", err)
	}

	readmeContent := GenerateReadme(client, mainHCL, variablesDefinitionHCL, resourceGetAPISpec)
	if err := writeContentToExamplesFolder(readmeContent, "README.md", resourceName); err != nil {
		log.Fatalf("Error writing main.tf: %v", err)
	}
}


func getOsArg() *string {
	if len(os.Args) < 2 {
		return nil
	}
	return &os.Args[1]
}

func GenerateVariableDefsHCL(client *openai.Client, mainHCLContent string) string {
	userPrompt := prompts.GetVarsDefHCLGenerationUserPrompt(prompts.VarsDefHCLUserPromptInputs{
		HCLConfig: mainHCLContent,
	})
	return CallModel(client, prompts.GenerateVarsDefHCLSystemPrompt, userPrompt)
}

func GenerateReadme(client *openai.Client, mainHCLContent, variablesDefinitionHCL, apiSpec string) string {
	userPrompt := prompts.GetReadmeGenerationUserPrompt(prompts.ReadmeUserPromptInputs{
		HCLConfig:       mainHCLContent,
		VariablesDefHCL: variablesDefinitionHCL,
		ResourceAPISpec: apiSpec,
	})
	return CallModel(client, prompts.GenerateReadmeSystemPrompt, userPrompt)
}

func GenerateMainHCL(client *openai.Client, resourceName string, resourceImpl string, apiSpec string) string {
	userPrompt := prompts.GetMainHCLGenerationUserPrompt(prompts.MainHCLUserPromptInputs{
		ResourceName:           resourceName,
		ResourceImplementation: resourceImpl,
		ResourceAPISpec:        apiSpec,
	})
	if err := os.WriteFile("hcl-generation-user-prompt.log", []byte(userPrompt), 0644); err != nil {
		log.Fatalf("failed to write to user prompt log file: %v", err)
	}
	return CallModel(client, prompts.GenerateMainHCLSystemPrompt, userPrompt)
}

func writeContentToExamplesFolder(content, fileName string, resourceName string) error {
	// Ensure the directory exists
	examplesDir := fmt.Sprintf("./examples/mongodbatlas_%s", resourceName)
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the HCL content to the file
	filePath := filepath.Join(examplesDir, fileName)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Printf("Content written to: %s\n", filePath)
	return nil
}

func getAPISpecSchema(resourceName string) string {
	getPath, ok := resourceToGetPath[resourceName]
	if !ok {
		return ""
	}

	if err := openapi.DownloadOpenAPISpec(atlasAdminAPISpecURL, specFilePath); err != nil {
		log.Fatalf("an error occurred when downloading Atlas Admin API spec: %v", err)
	}

	apiSpec, err := openapi.ParseAtlasAdminAPI(specFilePath)
	if err != nil {
		log.Fatalf("unable to parse Atlas Admin API: %v", err)
	}
	path, _ := apiSpec.Model.Paths.PathItems.Get(getPath)
	getOp := path.Get
	getOpBytes, _ := getOp.RenderInline()
	// okResponse, _ := getOp.Responses.Codes.Get(codespec.OASResponseCodeOK)
	// schema, _ := codespec.GetSchemaFromMediaType(okResponse.Content)
	// baseSchema := schema.Schema
	// schemaBytes, _ := baseSchema.RenderInline()
	return string(getOpBytes)
}

func getResourceImplementationSchema(resourceName string) string {
	lowerCaseJoined := codespec.SnakeCaseString(resourceName).LowerCaseNoUnderscore()
	var implementationContent string
	dirPath := resourcesBasePath + lowerCaseJoined
	files, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		log.Printf("Failed to read directory %s: %v", dirPath, err)
		return ""
	}
	for _, file := range files {
		baseName := filepath.Base(file)
		if strings.HasPrefix(baseName, "resource") {
			content, err := os.ReadFile(file)
			if err != nil {
				log.Printf("Failed to read file %s: %v", file, err)
				continue
			}

			implementationContent += fmt.Sprintf("File: %s\n\n%s\n\n", baseName, string(content))
		}
	}
	return implementationContent
}

const providerHCLContent = `provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}`

const versionsHCLContent = `terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.0"
    }
  }
  required_version = ">= 1.0"
}`
