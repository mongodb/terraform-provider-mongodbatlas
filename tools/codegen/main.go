package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
	"gopkg.in/yaml.v3"
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

		log.Printf("[INFO] Generating resource code: %s", resourceModel.Name)

		packageDir := fmt.Sprintf("internal/serviceapi/%s", resourceModel.PackageName)
		var generatedFiles []string

		schemaCode, err := schema.GenerateGoCode(resourceModel)
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
		schemaFilePath := fmt.Sprintf("%s/resource_schema.go", packageDir)
		if err := writeToFile(schemaFilePath, schemaCode); err != nil {
			log.Fatalf("[ERROR] An error occurred when writing content to file: %v", err)
		}
		generatedFiles = append(generatedFiles, schemaFilePath)

		resourceCode, err := resource.GenerateGoCode(resourceModel)
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
		resourceFilePath := fmt.Sprintf("%s/resource.go", packageDir)
		if err := writeToFile(resourceFilePath, resourceCode); err != nil {
			log.Fatalf("[ERROR] An error occurred when writing content to file: %v", err)
		}
		generatedFiles = append(generatedFiles, resourceFilePath)

		// Generate data source code if data sources are defined
		if resourceModel.DataSources != nil {
			log.Printf("[INFO] Generating data source code: %s", resourceModel.Name)

			// Generate data_source_schema.go
			dsSchemaCode, err := schema.GenerateDataSourceSchemaGoCode(resourceModel)
			if err != nil {
				log.Fatalf("[ERROR] %v", err)
			}
			dsSchemaFilePath := fmt.Sprintf("%s/data_source_schema.go", packageDir)
			if err := writeToFile(dsSchemaFilePath, dsSchemaCode); err != nil {
				log.Fatalf("[ERROR] An error occurred when writing content to file: %v", err)
			}
			generatedFiles = append(generatedFiles, dsSchemaFilePath)

			// Generate data_source.go
			dataSourceCode, err := datasource.GenerateGoCode(resourceModel)
			if err != nil {
				log.Fatalf("[ERROR] %v", err)
			}
			dataSourceFilePath := fmt.Sprintf("%s/data_source.go", packageDir)
			if err := writeToFile(dataSourceFilePath, dataSourceCode); err != nil {
				log.Fatalf("[ERROR] An error occurred when writing content to file: %v", err)
			}
			generatedFiles = append(generatedFiles, dataSourceFilePath)
		}

		// Format all generated files: goimports per file, fieldalignment on package
		formatGeneratedFiles(generatedFiles, packageDir)
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

// formatGeneratedFiles runs goimports on each file and fieldalignment on the package directory.
// Running fieldalignment on the package (not individual files) allows it to resolve types across files.
func formatGeneratedFiles(files []string, packageDir string) {
	// Run goimports on each file individually
	for _, filePath := range files {
		goimportsCmd := exec.CommandContext(context.Background(), "goimports", "-w", filePath)
		if output, err := goimportsCmd.CombinedOutput(); err != nil {
			log.Printf("[WARN] Goimports failed for %s: %v\nOutput: %s", filePath, err, output)
		}
	}

	// Run fieldalignment on the package directory so it can resolve types across files
	// Use filepath.Clean to sanitize the path and prevent command injection
	packagePath := "./" + packageDir
	fieldalignmentCmd := exec.CommandContext(context.Background(), "fieldalignment", "-fix", packagePath)
	if output, err := fieldalignmentCmd.CombinedOutput(); err != nil {
		log.Printf("[WARN] Fieldalignment failed for %s: %v\nOutput: %s", packagePath, err, output)
	}
}
