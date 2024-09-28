//nolint:gocritic
package codespec

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
)

func ConvertToProviderSpec(openAPIModel *openapi3.T, config config.Config, resourceName *string) *CodeSpecification {
	// var resourceSpec genconfig.Resource
	// if resourceName != nil {
	// 	resourceSpec = config.Resources[*resourceName]
	// }

	// TODO: convert openAPIModel and config to  CodeSpecification

	return &CodeSpecification{}
}
