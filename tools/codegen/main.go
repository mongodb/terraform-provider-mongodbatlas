package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
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
		schemaCode := schema.GenerateGoCode(&resourceModel, false) // object types are not needed as part of fully generated resources
		schemaFilePath := fmt.Sprintf("internal/serviceapi/%s/resource_schema.go", resourceModel.Name.LowerCaseNoUnderscore())
		if err := writeToFile(schemaFilePath, schemaCode); err != nil {
			log.Fatalf("an error occurred when writing content to file: %v", err)
		}
		formatGoFile(schemaFilePath)

		resourceCode := resource.GenerateGoCode(&resourceModel)
		resourceFilePath := fmt.Sprintf("internal/serviceapi/%s/resource.go", resourceModel.Name.LowerCaseNoUnderscore())
		if err := writeToFile(resourceFilePath, resourceCode); err != nil {
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
	// read/write/execute for owner, and read/execute for group and others
	const filePermission = 0o755

	// Create directories if they don't exist
	dir := filepath.Dir(fileName)
	dirPermission := os.FileMode(filePermission)
	if err := os.MkdirAll(dir, dirPermission); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write content to file (will override content if file exists)
	if err := os.WriteFile(fileName, []byte(content), filePermission); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", fileName, err)
	}
	return nil
}

// formatGoFile runs goimports and fieldalignment on the specified Go file
func formatGoFile(filePath string) {
	goimportsCmd := exec.CommandContext(context.Background(), "goimports", "-w", filePath)
	if output, err := goimportsCmd.CombinedOutput(); err != nil {
		log.Printf("warning: goimports failed for %s: %v\nOutput: %s", filePath, err, output)
	}

	fieldalignmentCmd := exec.CommandContext(context.Background(), "fieldalignment", "-fix", filePath)
	if output, err := fieldalignmentCmd.CombinedOutput(); err != nil {
		log.Printf("warning: fieldalignment failed for %s: %v\nOutput: %s", filePath, err, output)
	}
}
