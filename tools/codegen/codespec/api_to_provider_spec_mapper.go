//nolint:gocritic
package codespec

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
)

// using blank identifiers for now, will be removed in follow-up PRs once logic for conversion is added
func ToProviderSpecModel(atlasAdminAPISpecFilePath, configPath string, resourceName *string) *CodeSpecification {
	_, err := openapi.ParseAtlasAdminAPI(atlasAdminAPISpecFilePath)
	if err != nil {
		panic(err)
	}

	genConfig, _ := config.ParseGenConfigYAML(configPath)

	// var resourceSpec config.Resource
	if resourceName != nil {
		_ = genConfig.Resources[*resourceName]
	}

	return &CodeSpecification{}
}
