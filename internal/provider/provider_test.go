package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	providerfw "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
)

func TestResourceSchemas(t *testing.T) {
	t.Parallel()
	ctxProvider := context.Background()
	prov := provider.NewFrameworkProvider(nil)
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
			validateDocumentation(metadataRes.TypeName, schemaResponse)

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
	prov := provider.NewFrameworkProvider(nil)
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
			validateDSDocumentation(metadataRes.TypeName, schemaResponse)

			if schemaResponse.Diagnostics.HasError() {
				t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
			}

			if diagnostics := schemaResponse.Schema.ValidateImplementation(ctx); diagnostics.HasError() {
				t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
			}
		})
	}
}

func validateDocumentation(name string, resp *resource.SchemaResponse) {
	checkDescriptor(name, resp.Schema, &resp.Diagnostics)
	validateAttributes(name, resp.Schema.GetAttributes(), &resp.Diagnostics)
	validateBlocks(name, resp.Schema.GetBlocks(), &resp.Diagnostics)
}

func validateDSDocumentation(name string, resp *datasource.SchemaResponse) {
	checkDescriptor(name, resp.Schema, &resp.Diagnostics)
	validateAttributes(name, resp.Schema.GetAttributes(), &resp.Diagnostics)
	validateBlocks(name, resp.Schema.GetBlocks(), &resp.Diagnostics)
}

func validateAttribute(name string, attr schema.Attribute, diagnostics *diag.Diagnostics) {
	checkDescriptor(name, attr, diagnostics)
	if nested, ok := attr.(schema.NestedAttribute); ok {
		validateAttributes(name, nested.GetNestedObject().GetAttributes(), diagnostics)
	}
}

func validateBlock(name string, block schema.Block, diagnostics *diag.Diagnostics) {
	checkDescriptor(name, block, diagnostics)
	validateAttributes(name, block.GetNestedObject().GetAttributes(), diagnostics)
	validateBlocks(name, block.GetNestedObject().GetBlocks(), diagnostics)
}

func validateAttributes[T schema.Attribute](name string, attrs map[string]T, diagnostics *diag.Diagnostics) {
	for k, v := range attrs {
		validateAttribute(name+"."+k, v, diagnostics)
	}
}

func validateBlocks[T schema.Block](name string, blocks map[string]T, diagnostics *diag.Diagnostics) {
	for k, v := range blocks {
		validateBlock(name+"."+k, v, diagnostics)
	}
}

type descriptor interface {
	GetDescription() string
	GetMarkdownDescription() string
}

func checkDescriptor(name string, d descriptor, diagnostics *diag.Diagnostics) {
	if d.GetDescription() != d.GetMarkdownDescription() {
		diagnostics.Append(diag.NewErrorDiagnostic(
			"Conflicting Attribute Description",
			fmt.Sprintf("Description and MarkdownDescription differ for %q.", name),
		))
	}
}
