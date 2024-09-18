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

func validateDocumentation(resourceName string, schemaResponse *resource.SchemaResponse) {
	s := schemaResponse.Schema
	for attributeName, attribute := range s.GetAttributes() {
		validateAttribute(attribute, resourceName, attributeName, &schemaResponse.Diagnostics)
	}
	for blockName, block := range s.GetBlocks() {
		validateBlock(block, resourceName, blockName, &schemaResponse.Diagnostics)
	}
}

func validateDSDocumentation(resourceName string, schemaResponse *datasource.SchemaResponse) {
	s := schemaResponse.Schema
	for attributeName, attribute := range s.GetAttributes() {
		validateAttribute(attribute, resourceName, attributeName, &schemaResponse.Diagnostics)
	}
	for blockName, block := range s.GetBlocks() {
		validateBlock(block, resourceName, blockName, &schemaResponse.Diagnostics)
	}
}

func validateBlock(block schema.Block, resourceName, attributeName string, diagnostics *diag.Diagnostics) {
	if block.GetDescription() != block.GetMarkdownDescription() {
		diagnostics.Append(attributeIncorrectDescription(resourceName, attributeName))
	}
	for nestedAttributeName, nestedAttribute := range block.GetNestedObject().GetAttributes() {
		validateAttribute(nestedAttribute, resourceName, attributeName+"."+nestedAttributeName, diagnostics)
	}
	for nestedBlockName, nestedBlock := range block.GetNestedObject().GetBlocks() {
		validateBlock(nestedBlock, resourceName, attributeName+"."+nestedBlockName, diagnostics)
	}
}

func validateAttribute(attr schema.Attribute, resourceName, attributeName string, diagnostics *diag.Diagnostics) {
	if attr.GetDescription() != attr.GetMarkdownDescription() {
		diagnostics.Append(attributeIncorrectDescription(resourceName, attributeName))
	}
	if nested, ok := attr.(schema.NestedAttribute); ok {
		for nestedAttributeName, nestedAttribute := range nested.GetNestedObject().GetAttributes() {
			validateAttribute(nestedAttribute, resourceName, attributeName+"."+nestedAttributeName, diagnostics)
		}
	}
}

func attributeIncorrectDescription(resourceName, attributeName string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Incorrect Attribute Description",
		fmt.Sprintf("The Description and MarkdownDescription fields must be the same for %q.%q.", resourceName, attributeName),
	)
}
