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
	"gopkg.in/yaml.v3"
)

const (
	atlasAdminAPISpecURL        = "https://raw.githubusercontent.com/mongodb/atlas-sdk-go/main/openapi/atlas-api-transformed.yaml"
	configPath                  = "tools/codegen/config.yml"
	specFilePath                = "tools/codegen/open-api-spec.yml"
	resourceModelDir            = "tools/codegen/models/"
	resourceModelFilePathFormat = resourceModelDir + "%s.yaml"
)

func main() {
	resourceName := getOsArg()

	if err := openapi.DownloadOpenAPISpec(atlasAdminAPISpecURL, specFilePath); err != nil {
		log.Fatalf("[ERROR] An error occurred when downloading Atlas Admin API spec: %v", err)
	}

	{
		// Generate resource models
		model, err := codespec.ToCodeSpecModel(specFilePath, configPath, resourceName)
		if err != nil {
			log.Fatalf("[ERROR] An error occurred while generating codespec.Model: %v", err)
		}

		// Write resource models to files
		for i := range model.Resources {
			resourceModel := model.Resources[i]
			resourceModelFilePath := fmt.Sprintf(resourceModelFilePathFormat, resourceModel.Name)
			resourceModelYaml, err := yaml.Marshal(resourceModel)
			if err != nil {
				log.Fatalf("[ERROR] An error occurred while serializing the resource model: %v", err)
			}
			if err := writeToFile(resourceModelFilePath, resourceModelYaml); err != nil {
				log.Fatalf("[ERROR] An error occurred while writing resource model to file: %v", err)
			}
		}
	}

	// Gather resource model files
	var resourceModelFilePaths []string
	if resourceName == nil {
		files, err := os.ReadDir(resourceModelDir)
		if err != nil {
			log.Fatalf("[ERROR] An error occurred while reading the resource model directory: %v", err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			resourceModelFilePaths = append(resourceModelFilePaths, resourceModelDir+file.Name())
		}
	} else {
		resourceModelFilePaths = append(resourceModelFilePaths, fmt.Sprintf(resourceModelFilePathFormat, *resourceName))
	}

	// Generate code from resource model files
	for _, filePath := range resourceModelFilePaths {
		resourceModel, err := readResourceModelFromFile(filePath)
		if err != nil {
			log.Fatalf("[ERROR] An error occurred while reading the resource model file: %v", err)
		}

		log.Printf("[INFO] Generating resource code: %s", resourceModel.Name)

		schemaCode := schema.GenerateGoCode(resourceModel, false) // object types are not needed as part of fully generated resources
		schemaFilePath := fmt.Sprintf("internal/serviceapi/%s/resource_schema.go", resourceModel.Name.LowerCaseNoUnderscore())
		if err := writeToFile(schemaFilePath, schemaCode); err != nil {
			log.Fatalf("[ERROR] An error occurred when writing content to file: %v", err)
		}
		formatGoFile(schemaFilePath)

		resourceCode := resource.GenerateGoCode(resourceModel)
		resourceFilePath := fmt.Sprintf("internal/serviceapi/%s/resource.go", resourceModel.Name.LowerCaseNoUnderscore())
		if err := writeToFile(resourceFilePath, resourceCode); err != nil {
			log.Fatalf("[ERROR] An error occurred when writing content to file: %v", err)
		}
	}
}

func getOsArg() *string {
	if len(os.Args) < 2 {
		return nil
	}
	return &os.Args[1]
}

func writeToFile(filePath string, content []byte) error {
	// read/write/execute for owner, and read/execute for group and others
	const filePermission = 0o755

	// Create directories if they don't exist
	dir := filepath.Dir(filePath)
	dirPermission := os.FileMode(filePermission)
	if err := os.MkdirAll(dir, dirPermission); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write content to file (will override content if file exists)
	if err := os.WriteFile(filePath, content, filePermission); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	return nil
}

func readResourceModelFromFile(filePath string) (*codespec.Resource, error) {
	resourceModelYaml, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var resourceModel codespec.Resource
	err = yaml.Unmarshal(resourceModelYaml, &resourceModel)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize resource: %w", err)
	}
	return &resourceModel, nil
}

// formatGoFile runs goimports and fieldalignment on the specified Go file
func formatGoFile(filePath string) {
	goimportsCmd := exec.CommandContext(context.Background(), "goimports", "-w", filePath)
	if output, err := goimportsCmd.CombinedOutput(); err != nil {
		log.Printf("[WARN] Goimports failed for %s: %v\nOutput: %s", filePath, err, output)
	}

	fieldalignmentCmd := exec.CommandContext(context.Background(), "fieldalignment", "-fix", filePath)
	if output, err := fieldalignmentCmd.CombinedOutput(); err != nil {
		log.Printf("[WARN] Fieldalignment failed for %s: %v\nOutput: %s", filePath, err, output)
	}
}
