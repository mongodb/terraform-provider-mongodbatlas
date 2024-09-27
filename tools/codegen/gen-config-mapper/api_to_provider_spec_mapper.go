//nolint:gocritic
package genconfigmapper

import (
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	genconfig "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gen-config"
)

func ConvertToProviderSpec(openAPIModel *libopenapi.DocumentModel[v3.Document], config genconfig.Config, resourceName *string) *CodeSpecification {
	// var resourceSpec genconfig.Resource
	// if resourceName != nil {
	// 	resourceSpec = config.Resources[*resourceName]
	// }

	// TODO: convert openAPIModel and config to  CodeSpecification

	return &CodeSpecification{}
}
