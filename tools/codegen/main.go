package main

import (
	"fmt"
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
		fmt.Println("No resource name provided")
	} else {
		fmt.Printf("Resource name: %s\n", *resourceName)
	}

	if err := openapi.DownloadOpenAPISpec(atlasAdminAPISpecURL, specFilePath); err != nil {
		panic(err)
	}

	// apiDocModel, err := openapi.ParseAtlasAdminAPI(atlasAdminAPISpecURL)
	// if err != nil {
	// 	panic(err)
	// }

	// genConfig, _ := config.ParseGenConfigYAML(configPath)

	_ = codespec.ToProviderSpec(specFilePath, configPath, resourceName)
}

func getOsArg() *string {
	if len(os.Args) < 2 {
		return nil
	}
	return &os.Args[1]
}
