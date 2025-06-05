package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
)

var (
	configPath  = "tools/codegen/config.yml"
	specDirPath = "tools/codegen/open-api-specs"
)

func main() {
	resourceName := getOsArg()
	envConfigPath := os.Getenv("CODEGEN_CONFIG_PATH")
	if envConfigPath != "" {
		log.Printf("Using custom codegen config file path: %s", envConfigPath)
		configPath = envConfigPath
	}
	configModel, err := config.ParseGenConfigYAML(configPath)
	if err != nil {
		log.Fatalf("unable to parse config file: %v", err)
	}
	envSpecDirPath := os.Getenv("OPENAPI_SPEC_DIR_PATH")
	if envSpecDirPath != "" {
		specDirPath = envSpecDirPath
		log.Printf("Using custom OpenAPI spec dir path: %s", specDirPath)
	}
	skipOpenAPIDownload, _ := strconv.ParseBool(os.Getenv("SKIP_OPENAPI_DOWNLOAD"))
	if skipOpenAPIDownload {
		log.Println("Skipping download of Atlas Admin API spec")
	} else if err := openapi.DownloadOpenAPISpecs(configModel.APISpecs, specDirPath); err != nil {
		log.Fatalf("an error occurred when downloading Atlas Admin API spec: %v", err)
	}
	apiSpecsParsed := openapi.ParseAPISpecs(specDirPath, configModel.APISpecsNames())
	model, err := codespec.ToCodeSpecModel(apiSpecsParsed, configModel, resourceName)
	if err != nil {
		log.Fatalf("an error occurred while generating codespec.Model: %v", err)
	}

	for i := range model.Resources {
		resourceModel := model.Resources[i]
		schemaCode := schema.GenerateGoCode(&resourceModel, false) // object types are not needed as part of fully generated resources
		if err := writeToFile(fmt.Sprintf("internal/serviceapi/%s/resource_schema.go", resourceModel.Name.LowerCaseNoUnderscore()), schemaCode); err != nil {
			log.Fatalf("an error occurred when writing content to file: %v", err)
		}
		resourceCode := resource.GenerateGoCode(&resourceModel)
		if err := writeToFile(fmt.Sprintf("internal/serviceapi/%s/resource.go", resourceModel.Name.LowerCaseNoUnderscore()), resourceCode); err != nil {
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
