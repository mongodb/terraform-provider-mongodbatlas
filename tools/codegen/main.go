package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/schema"
	"gopkg.in/yaml.v3"
)

const (
	atlasAdminAPISpecURL = "https://raw.githubusercontent.com/mongodb/atlas-sdk-go/main/openapi/atlas-api-transformed.yaml"
	configPath           = "tools/codegen/config.yml"
	specFilePath         = "tools/codegen/open-api-spec.yml"
)

func main() {
	resourceName := getOsArg()
	skipOpenAPIDownload := os.Getenv("SKIP_OPENAPI_DOWNLOAD")
	if skipOpenAPIDownload == "true" {
		log.Println("Skipping download of Atlas Admin API spec")
	} else if err := openapi.DownloadOpenAPISpec(atlasAdminAPISpecURL, specFilePath); err != nil {
		log.Fatalf("an error occurred when downloading Atlas Admin API spec: %v", err)
	}

	model, err := codespec.ToCodeSpecModel(specFilePath, configPath, resourceName)
	if err != nil {
		log.Fatalf("an error occurred while generating codespec.Model: %v", err)
	}

	for i := range model.Resources {
		resourceModel := model.Resources[i]
		err := dumpYaml(resourceModel, resourceModel.Name.LowerCaseNoUnderscore()+".yml")
		if err != nil {
			log.Fatalf("an error occurred when dumping yaml: %v", err)
		}
		schemaCode := schema.GenerateGoCode(resourceModel)
		if err := writeToFile(fmt.Sprintf("internal/service/%s/resource_schema.go", resourceModel.Name.LowerCaseNoUnderscore()), schemaCode); err != nil {
			log.Fatalf("an error occurred when writing content to file: %v", err)
		}
	}
}

func dumpYaml(myStruct any, filename string) error {
	dirPath := os.Getenv("SPEC_OUTPUT_DIR")
	filePath := dirPath + filename
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		log.Fatalf("error opening/creating file: %v", err)
	}
	defer file.Close()
	enc := yaml.NewEncoder(file)
	err = enc.Encode(myStruct)
	if err != nil {
		return err
	}
	yamlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fmt.Printf("marshaled struct %s", string(yamlBytes))
	return err
}

func getOsArg() *string {
	if len(os.Args) < 2 {
		return nil
	}
	return &os.Args[1]
}

func writeToFile(fileName, content string) error {
	// will override content if file exists
	err := os.WriteFile(fileName, []byte(content), 0o600)
	if err != nil {
		return err
	}
	return nil
}
