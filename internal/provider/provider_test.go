package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	providerfw "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
)

func TestResourceSchemas(t *testing.T) {
	t.Parallel()
	ctxProvider := context.Background()
	prov := provider.NewFrameworkProvider()
	var provReq providerfw.MetadataRequest
	var provRes providerfw.MetadataResponse
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

func TestDataSourceSchemas(t *testing.T) {
	t.Parallel()
	ctxProvider := context.Background()
	prov := provider.NewFrameworkProvider()
	var provReq providerfw.MetadataRequest
	var provRes providerfw.MetadataResponse
	prov.Metadata(ctxProvider, provReq, &provRes)
	for _, fn := range prov.DataSources(ctxProvider) {
		ctx := context.Background()
		res := fn()
		metadataReq := datasource.MetadataRequest{
			ProviderTypeName: provRes.TypeName,
		}
		var metadataRes datasource.MetadataResponse
		res.Metadata(ctx, metadataReq, &metadataRes)

		t.Run(metadataRes.TypeName, func(t *testing.T) {
			schemaRequest := datasource.SchemaRequest{}
			schemaResponse := &datasource.SchemaResponse{}
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
