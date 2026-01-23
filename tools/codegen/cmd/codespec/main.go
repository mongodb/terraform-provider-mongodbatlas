package main

import (
	"flag"
	"fmt"
	"log"

	"go.yaml.in/yaml/v4"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	resourceName, resourceTier, err := getArgs()
	if err != nil {
		log.Fatalf("[ERROR] Invalid resource tier: %v", err)
	}
	if resourceName != nil {
		log.Printf("Generating resource models for resource: %s", *resourceName)
	}
	if resourceTier != nil {
		log.Printf("Using resource tier filter: %s", *resourceTier)
	}
	if err := writeResourceModels(resourceName, resourceTier); err != nil {
		log.Fatalf("[ERROR] An error occurred while generating resource models: %v", err)
	}
}

func getArgs() (resourceName *string, resourceTier *codespec.ResourceTier, err error) {
	var resourceNameFlag string
	var resourceTierFlag string
	flag.StringVar(&resourceNameFlag, "resource-name", "", "Generate models only for the specified resource name")
	flag.StringVar(&resourceTierFlag, "resource-tier", "", "Generate models only for resources in the specified tier (prod|internal)")
	flag.Parse()
	resourceTier, err = codespec.ParseResourceTier(resourceTierFlag)
	if err != nil {
		return nil, nil, err
	}
	return conversion.StringPtr(resourceNameFlag), resourceTier, nil
}

func writeResourceModels(resourceName *string, resourceTier *codespec.ResourceTier) error {
	// Generate resource and data source models from API spec
	model, err := codespec.ToCodeSpecModel(SpecFilePath, ConfigPath, resourceName, resourceTier)
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
