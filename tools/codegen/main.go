package main

import (
	"fmt"
	"os"
)

const (
	atlasAdminAPISpecURL = "https://raw.githubusercontent.com/mongodb/atlas-sdk-go/main/openapi/atlas-api-transformed.yaml"
	configPath           = "schema-generation/config.yml"
)

func main() {
	resourceName := getOsArg()
	if resourceName == nil {
		fmt.Println("No resource name provided")
	} else {
		fmt.Printf("Resource name: %s\n", *resourceName)
	}

	// apiDocModel, err := openapi.ParseAtlasAdminAPI(atlasAdminAPISpecURL)
	// if err != nil {
	// 	panic(err)
	// }

	// genConfig, _ := config.ParseGenConfigYAML(configPath)

	// _ = codespec.ConvertToProviderSpec(apiDocModel, *genConfig, resourceName)
}

func getOsArg() *string {
	if len(os.Args) < 2 {
		return nil
	}
	return &os.Args[1]
}
