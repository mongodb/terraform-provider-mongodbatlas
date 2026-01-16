package main

import (
	"fmt"
	"log"
	"os"

	"go.yaml.in/yaml/v4"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/fileutil"
)

const (
	ConfigPath                  = "tools/codegen/config.yml"
	SpecFilePath                = "tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml"
	ResourceModelDir            = "tools/codegen/models/"
	ResourceModelFilePathFormat = ResourceModelDir + "%s.yaml"
)

func main() {
	resourceName := getResourceNameArg()
	if err := writeResourceModels(resourceName); err != nil {
		log.Fatalf("[ERROR] An error occurred while generating resource models: %v", err)
	}
}

func getResourceNameArg() *string {
	if len(os.Args) < 2 {
		return nil
	}
	return &os.Args[1]
}

func writeResourceModels(resourceName *string) error {
	// Generate resource and data source models from API spec
	model, err := codespec.ToCodeSpecModel(SpecFilePath, ConfigPath, resourceName)
	if err != nil {
		return fmt.Errorf("failed to generate codespec.Model: %w", err)
	}

	// Write resource models to files
	for i := range model.Resources {
		resourceModel := model.Resources[i]
		resourceModelFilePath := fmt.Sprintf(ResourceModelFilePathFormat, resourceModel.Name)
		resourceModelYaml, err := yaml.Marshal(resourceModel)
		if err != nil {
			return fmt.Errorf("failed to serialize resource model: %w", err)
		}
		if err := fileutil.WriteFile(resourceModelFilePath, resourceModelYaml); err != nil {
			return fmt.Errorf("failed to write resource model to file: %w", err)
		}
	}

	return nil
}
