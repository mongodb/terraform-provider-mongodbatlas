package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen"
)

const (
	ConfigPath                  = "tools/codegen/config.yml"
	SpecFilePath                = "tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml"
	ResourceModelDir            = "tools/codegen/models/"
	ResourceModelFilePathFormat = ResourceModelDir + "%s.yaml"

	StepModelGen = "model-gen"
	StepCodeGen  = "code-gen"
)

func main() {
	resourceName, resourceTier, step, err := getArgs()
	if err != nil {
		log.Fatalf("[ERROR] Invalid arguments: %v", err)
	}
	if resourceName != nil {
		log.Printf("Resource filter: %s", *resourceName)
	}
	if resourceTier != nil {
		log.Printf("Resource tier filter: %s", *resourceTier)
	}
	if step != nil {
		log.Printf("Executing individual step: %s", *step)
	}

	if err := runSteps(step, resourceName, resourceTier); err != nil {
		log.Fatalf("[ERROR] Code generation failed: %v", err)
	}
}

func getArgs() (resourceName *string, resourceTier *codespec.ResourceTier, step *string, err error) {
	var resourceNameFlag string
	var resourceTierFlag string
	var stepFlag string
	flag.StringVar(&resourceNameFlag, "resource-name", "", "Generate models only for the specified resource name")
	flag.StringVar(&resourceTierFlag, "resource-tier", "", "Generate models only for resources in the specified tier (prod|internal)")
	flag.StringVar(&stepFlag, "step", "", "Run only one step: model-gen or code-gen (default runs both)")
	flag.Parse()

	step, err = parseStep(stepFlag)
	if err != nil {
		return nil, nil, nil, err
	}
	resourceTier, err = codespec.ParseResourceTier(resourceTierFlag)
	if err != nil {
		return nil, nil, nil, err
	}
	return conversion.StringPtr(resourceNameFlag), resourceTier, step, nil
}

func parseStep(step string) (*string, error) {
	if step == "" {
		return nil, nil
	}
	switch step {
	case StepModelGen, StepCodeGen:
		return &step, nil
	default:
		return nil, fmt.Errorf("invalid step %q (allowed: %s, %s)", step, StepModelGen, StepCodeGen)
	}
}

func runSteps(step, resourceName *string, resourceTier *codespec.ResourceTier) error {
	if step == nil || *step == StepModelGen {
		if err := writeResourceModels(resourceName, resourceTier); err != nil {
			return err
		}
	}
	if step == nil || *step == StepCodeGen {
		return generateFromResourceModels(resourceName, resourceTier)
	}
	return nil
}

func writeResourceModels(resourceName *string, resourceTier *codespec.ResourceTier) error {
	model, err := codespec.ToCodeSpecModel(SpecFilePath, ConfigPath, resourceName, resourceTier)
	if err != nil {
		return fmt.Errorf("failed to generate codespec.Model: %w", err)
	}

	for i := range model.Resources {
		resourceModel := model.Resources[i]
		resourceModelFilePath := fmt.Sprintf(ResourceModelFilePathFormat, resourceModel.Name)
		resourceModelYaml, err := yaml.Marshal(resourceModel)
		if err != nil {
			return fmt.Errorf("failed to serialize resource model: %w", err)
		}
		if err := writeFile(resourceModelFilePath, resourceModelYaml); err != nil {
			return fmt.Errorf("failed to write resource model to file: %w", err)
		}
	}

	return nil
}

func generateFromResourceModels(resourceName *string, resourceTier *codespec.ResourceTier) error {
	resourceModelFilePaths, err := gatherResourceModelFilePaths(resourceName, resourceTier)
	if err != nil {
		return fmt.Errorf("failed to gather resource model files: %w", err)
	}
	if err := generateCodeFromResourceModels(resourceModelFilePaths); err != nil {
		return fmt.Errorf("failed to generate code from resource models: %w", err)
	}
	return nil
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
	for _, filePath := range resourceModelFilePaths {
		resourceModel, err := readResourceModelFromFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read resource model file: %w", err)
		}

		packageDir := fmt.Sprintf("internal/serviceapi/%s", resourceModel.PackageName)

		if _, err := gofilegen.GenerateCodeForResource(resourceModel, packageDir, writeFile); err != nil {
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

func writeFile(filePath string, content []byte) error {
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
