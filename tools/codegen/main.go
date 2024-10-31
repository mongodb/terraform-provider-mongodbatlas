package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

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

	specDirPath := os.Getenv("SPEC_RESOURCE_OUTPUT_DIR")
	for i := range model.Resources {
		resourceModel := model.Resources[i]
		if specDirPath != "" {
			dumpYaml(resourceModel, path.Join(specDirPath, resourceModel.Name.LowerCaseNoUnderscore()+".yaml"))
		}
		if err != nil {
			log.Fatalf("an error occurred when dumping yaml: %v", err)
		}
		schemaCode := schema.GenerateGoCode(resourceModel)
		if err := writeToFile(fmt.Sprintf("internal/service/%s/resource_schema.go", resourceModel.Name.LowerCaseNoUnderscore()), schemaCode); err != nil {
			log.Fatalf("an error occurred when writing content to file: %v", err)
		}
	}
}

func dumpYaml(resource any, filePath string) error {
	initialYaml := strings.Builder{}
	// file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	// if err != nil {
	// 	log.Fatalf("error opening/creating file: %v", err)
	// }
	// defer file.Close()
	enc := yaml.NewEncoder(&initialYaml)
	err := enc.Encode(resource)
	if err != nil {
		return err
	}
	// yamlBytes, err := os.ReadFile(filePath)
	// if err != nil {
	// 	return err
	// }
	// yamlString := string(yamlBytes)
	yamlString := initialYaml.String()
	yamlStringCleaned := strings.Builder{}
	for _, line := range strings.Split(yamlString, "\n") {
		if strings.HasSuffix(line, ": null") {
			continue
		}
		yamlStringCleaned.WriteString(line + "\n")
	}
	stemName := path.Base(filePath)
	yamlFinal := yamlStringCleaned.String()
	if err != nil {
		return err
	}
	err = writeToFile(filePath, yamlFinal)
	// err = os.WriteFile(filePath, []byte(yamlFinal), 0o600)
	// _, err = file.WriteString(yamlFinal)
	if err != nil {
		return err
	}
	fmt.Printf("dumped resource %s to %s:\n%s", stemName, filePath, yamlFinal)
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
