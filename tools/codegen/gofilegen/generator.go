package gofilegen

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
)

// CodeGenerator defines a function that generates Go code from a resource model
type CodeGenerator func(*codespec.Resource) ([]byte, error)

// FileWriter defines a function that writes content to a file
type FileWriter func(filePath string, content []byte) error

// ComponentSpec defines a component (resource/data source) with its schema and implementation generators
type ComponentSpec struct {
	SchemaGenerator CodeGenerator // generates the _schema.go file
	ImplGenerator   CodeGenerator // generates the .go file
	BaseName        string        // one of "resource", "data_source", "plural_data_source"
}

var (
	resourceComponent = ComponentSpec{
		SchemaGenerator: schema.GenerateGoCode,
		ImplGenerator:   resource.GenerateGoCode,
		BaseName:        "resource",
	}

	dataSourceComponent = ComponentSpec{
		SchemaGenerator: schema.GenerateDataSourceSchemaGoCode,
		ImplGenerator:   datasource.GenerateGoCode,
		BaseName:        "data_source",
	}

	pluralDataSourceComponent = ComponentSpec{
		SchemaGenerator: schema.GeneratePluralDataSourceSchemaGoCode,
		ImplGenerator:   datasource.GeneratePluralGoCode,
		BaseName:        "plural_data_source",
	}
)

// GenerateCodeForResource generates all Go files for a resource and its data sources
func GenerateCodeForResource(resourceModel *codespec.Resource, packageDir string, writeFile FileWriter) ([]string, error) {
	log.Printf("[INFO] Generating resource code: %s", resourceModel.Name)

	var generatedFiles []string

	hasResourceOps := hasResourceOps(resourceModel.Operations)

	if hasResourceOps {
		// Generate resource files only if resource operations are defined
		files, err := generateComponentFiles(resourceModel, packageDir, resourceComponent, writeFile)
		if err != nil {
			return nil, fmt.Errorf("failed to generate resource: %w", err)
		}
		generatedFiles = append(generatedFiles, files...)
	}

	// Generate data source files if data sources are defined
	if resourceModel.DataSources != nil {
		log.Printf("[INFO] Generating data source code: %s", resourceModel.Name)

		files, err := generateComponentFiles(resourceModel, packageDir, dataSourceComponent, writeFile)
		if err != nil {
			return nil, fmt.Errorf("failed to generate data source: %w", err)
		}
		generatedFiles = append(generatedFiles, files...)

		// Generate plural data source files if plural data source is defined
		if isPluralDataSourceDefined(*resourceModel.DataSources) {
			files, err = generateComponentFiles(resourceModel, packageDir, pluralDataSourceComponent, writeFile)
			if err != nil {
				return nil, fmt.Errorf("failed to generate plural data source: %w", err)
			}
			generatedFiles = append(generatedFiles, files...)
		}
	}

	// Format all generated files: goimports per file, fieldalignment on package
	FormatGeneratedFiles(generatedFiles, packageDir)

	return generatedFiles, nil
}

func isPluralDataSourceDefined(dataSources codespec.DataSources) bool {
	return dataSources.Schema != nil && dataSources.Schema.PluralDSAttributes != nil && len(*dataSources.Schema.PluralDSAttributes) > 0
}

// generateComponentFiles generates schema and implementation files for a resource or data source
func generateComponentFiles(resourceModel *codespec.Resource, packageDir string, spec ComponentSpec, writeFile FileWriter) ([]string, error) {
	var generatedFiles []string

	schemaFile := fmt.Sprintf("%s_schema.go", spec.BaseName)
	schemaPath, err := generateAndWriteFile(resourceModel, packageDir, schemaFile, spec.SchemaGenerator, spec.BaseName+" schema", writeFile)
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, schemaPath)

	implFile := fmt.Sprintf("%s.go", spec.BaseName)
	implPath, err := generateAndWriteFile(resourceModel, packageDir, implFile, spec.ImplGenerator, spec.BaseName, writeFile)
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, implPath)

	return generatedFiles, nil
}

// generateAndWriteFile generates code using the provided generator and writes it to a file
func generateAndWriteFile(resourceModel *codespec.Resource, packageDir, fileName string, generator CodeGenerator, logMsg string, writeFile FileWriter) (string, error) {
	code, err := generator(resourceModel)
	if err != nil {
		return "", fmt.Errorf("failed to generate %s: %w", logMsg, err)
	}

	filePath := fmt.Sprintf("%s/%s", packageDir, fileName)
	if err := writeFile(filePath, code); err != nil {
		return "", fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return filePath, nil
}

// FormatGeneratedFiles runs goimports on each file and fieldalignment on the package directory.
// Running fieldalignment on the package (not individual files) allows it to resolve types across files.
func FormatGeneratedFiles(files []string, packageDir string) {
	// Run goimports on each file individually
	for _, filePath := range files {
		goimportsCmd := exec.CommandContext(context.Background(), "goimports", "-w", filePath)
		if output, err := goimportsCmd.CombinedOutput(); err != nil {
			log.Printf("[WARN] Goimports failed for %s: %v\nOutput: %s", filePath, err, output)
		}
	}

	// Run fieldalignment on the package directory so it can resolve types across files
	packagePath := "./" + packageDir
	fieldalignmentCmd := exec.CommandContext(context.Background(), "fieldalignment", "-fix", packagePath)
	if output, err := fieldalignmentCmd.CombinedOutput(); err != nil {
		log.Printf("[WARN] Fieldalignment failed for %s: %v\nOutput: %s", packagePath, err, output)
	}
}

// Checks if any resource operations are defined
func hasResourceOps(ops codespec.APIOperations) bool {
	return ops.Create != nil || ops.Read != nil || ops.Update != nil || ops.Delete != nil
}
