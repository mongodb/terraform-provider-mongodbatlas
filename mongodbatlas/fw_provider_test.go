package mongodbatlas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestResourceSchemas(t *testing.T) {
	t.Parallel()
	ctxProvider := context.Background()
	prov := NewFrameworkProvider()
	var provReq provider.MetadataRequest
	var provRes provider.MetadataResponse
	prov.Metadata(ctxProvider, provReq, &provRes)
	for _, fn := range prov.Resources(ctxProvider) {
		ctx := context.Background()
		res := fn()
		metadataReq := resource.MetadataRequest{
			ProviderTypeName: provRes.TypeName,
		}
		var metadataRes resource.MetadataResponse
		res.Metadata(ctx, metadataReq, &metadataRes)

		t.Run(metadataRes.TypeName, func(t *testing.T) {
			schemaRequest := resource.SchemaRequest{}
			schemaResponse := &resource.SchemaResponse{}

			res.Schema(ctx, schemaRequest, schemaResponse)

			if schemaResponse.Diagnostics.HasError() {
				t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
			}

			if diagnostics := schemaResponse.Schema.ValidateImplementation(ctx); diagnostics.HasError() {
				t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
			}
		})
	}
}
