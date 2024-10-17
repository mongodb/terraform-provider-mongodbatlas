package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/schema"
)

const (
	atlasAdminAPISpecURL = "https://raw.githubusercontent.com/mongodb/atlas-sdk-go/main/openapi/atlas-api-transformed.yaml"
	configPath           = "tools/codegen/config.yml"
	specFilePath         = "tools/codegen/open-api-spec.yml"
)

func main() {
	resourceName := getOsArg()

	if err := openapi.DownloadOpenAPISpec(atlasAdminAPISpecURL, specFilePath); err != nil {
		log.Fatalf("an error occurred when downloading Atlas Admin API spec: %v", err)
	}

	model, err := codespec.ToCodeSpecModel(specFilePath, configPath, resourceName)
	if err != nil {
		log.Fatalf("an error occurred while generating codespec.Model: %v", err)
	}

	for i := range model.Resources {
		resourceModel := model.Resources[i]
		schemaCode := schema.GenerateGoCode(resourceModel)
		if err := writeToFile(fmt.Sprintf("internal/service/%s/resource_schema.go", resourceModel.Name.LowerCaseNoUnderscore()), schemaCode); err != nil {
			log.Fatalf("an error occurred when writing content to file: %v", err)
		}
	}
}

func getOsArg() *string {
	if len(os.Args) < 2 {
		return nil
	}
	return &os.Args[1]
}

func writeToFile(fileName, content string) error {
	// Open the file with write-only permission, create it if it doesn't exist and truncate its content if it exists
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
