package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/fileutil"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen"
	"github.com/stretchr/testify/assert/yaml"
)

const (
	ResourceModelDir            = "tools/codegen/models/"
	ResourceModelFilePathFormat = ResourceModelDir + "%s.yaml"
)

func main() {
	resourceName := getResourceNameArg()
	resourceModelFilePaths, err := gatherResourceModelFilePaths(resourceName)
	if err != nil {
		log.Fatalf("[ERROR] An error occurred while gathering resource model files: %v", err)
	}
	if err := generateCodeFromResourceModels(resourceModelFilePaths); err != nil {
		log.Fatalf("[ERROR] An error occurred while generating code from resource models: %v", err)
	}
}

func getResourceNameArg() *string {
	if len(os.Args) < 2 {
		return nil
	}
	return &os.Args[1]
}

func gatherResourceModelFilePaths(resourceName *string) ([]string, error) {
	var resourceModelFilePaths []string
	if resourceName == nil {
		files, err := os.ReadDir(ResourceModelDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read resource model directory: %w", err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			resourceModelFilePaths = append(resourceModelFilePaths, ResourceModelDir+file.Name())
		}
		return resourceModelFilePaths, nil
	}

	resourceModelFilePaths = append(resourceModelFilePaths, fmt.Sprintf(ResourceModelFilePathFormat, *resourceName))
	return resourceModelFilePaths, nil
}

func generateCodeFromResourceModels(resourceModelFilePaths []string) error {
	// Generate code from resource model files
	for _, filePath := range resourceModelFilePaths {
		resourceModel, err := readResourceModelFromFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read resource model file: %w", err)
		}

		packageDir := fmt.Sprintf("internal/serviceapi/%s", resourceModel.PackageName)

		// Generate all files for the resource and its data sources
		if _, err := gofilegen.GenerateCodeForResource(resourceModel, packageDir, fileutil.WriteFile); err != nil {
			return fmt.Errorf("failed to generate code for %s: %w", resourceModel.Name, err)
		}
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
