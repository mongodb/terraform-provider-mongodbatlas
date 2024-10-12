package main

import (
	"log"
	"os"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
)

const (
	atlasAdminAPISpecURL = "https://raw.githubusercontent.com/mongodb/atlas-sdk-go/main/openapi/atlas-api-transformed.yaml"
	configPath           = "schema-generation/config.yml"
	specFilePath         = "open-api-spec.yml"
)

func main() {
	resourceName := getOsArg()
	if resourceName == nil {
		log.Fatal("No resource name provided")
	}
	log.Printf("Resource name: %s\n", *resourceName)

	if err := openapi.DownloadOpenAPISpec(atlasAdminAPISpecURL, specFilePath); err != nil {
		log.Fatalf("an error occurred when downloading Atlas Admin API spec: %v", err)
	}

	_, err := codespec.ToCodeSpecModel(specFilePath, configPath, *resourceName)
	if err != nil {
		log.Fatalf("an error occurred while generating codespec.Model: %v", err)
	}
}

func getOsArg() *string {
	if len(os.Args) < 2 {
		return nil
	}
	return &os.Args[1]
}
