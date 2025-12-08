package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen"
)

const (
	configPath                  = "tools/codegen/config.yml"
	specFilePath                = "tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml"
	resourceModelDir            = "tools/codegen/models/"
	resourceModelFilePathFormat = resourceModelDir + "%s.yaml"
)

func main() {
	resourceName := getOsArg()

	// Generate resource and data source models from API spec
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

		packageDir := fmt.Sprintf("internal/serviceapi/%s", resourceModel.PackageName)

		// Generate all files for the resource and its data sources
		if _, err := gofilegen.GenerateCodeForResource(resourceModel, packageDir, writeToFile); err != nil {
			log.Fatalf("[ERROR] Failed to generate code for %s: %v", resourceModel.Name, err)
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
