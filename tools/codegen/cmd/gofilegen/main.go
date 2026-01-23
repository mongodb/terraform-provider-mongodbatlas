package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	resourceName, resourceTier, err := getArgs()
	if err != nil {
		log.Fatalf("[ERROR] Invalid arguments: %v", err)
	}
	if resourceName != nil {
		log.Printf("Generating code for resource: %s", *resourceName)
	}
	resourceModelFilePaths, err := gatherResourceModelFilePaths(resourceName, resourceTier)
	if err != nil {
		log.Fatalf("[ERROR] An error occurred while gathering resource model files: %v", err)
	}
	if err := generateCodeFromResourceModels(resourceModelFilePaths); err != nil {
		log.Fatalf("[ERROR] An error occurred while generating code from resource models: %v", err)
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

func gatherResourceModelFilePaths(resourceName *string, resourceTier *codespec.ResourceTier) ([]string, error) {
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
	} else {
		resourceModelFilePaths = append(resourceModelFilePaths, fmt.Sprintf(ResourceModelFilePathFormat, *resourceName))
	}

	if resourceTier != nil {
		var filtered []string
		for _, filePath := range resourceModelFilePaths {
			isInternal := strings.HasSuffix(filePath, codespec.InternalResourceSuffix+".yaml")
			switch *resourceTier {
			case codespec.ResourceTierInternal:
				if isInternal {
					filtered = append(filtered, filePath)
				}
			case codespec.ResourceTierProd:
				if !isInternal {
					filtered = append(filtered, filePath)
				}
			}
		}
		resourceModelFilePaths = filtered
	}
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
