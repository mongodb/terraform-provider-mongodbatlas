package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type AttributeSimplified interface {
	// Implementations should include the tftypes.AttributePathStepper
	// interface methods for proper path and data handling.
	tftypes.AttributePathStepper

	// Equal should return true if the other attribute is exactly equivalent.
	Equal(o AttributeSimplified) bool

	// GetDeprecationMessage should return a non-empty string if an attribute
	// is deprecated. This is named differently than DeprecationMessage to
	// prevent a conflict with the tfsdk.Attribute field name.
	GetDeprecationMessage() string

	// GetDescription should return a non-empty string if an attribute
	// has a plaintext description. This is named differently than Description
	// to prevent a conflict with the tfsdk.Attribute field name.
	GetDescription() string

	// GetMarkdownDescription should return a non-empty string if an attribute
	// has a Markdown description. This is named differently than
	// MarkdownDescription to prevent a conflict with the tfsdk.Attribute field
	// name.
	GetMarkdownDescription() string

	// GetType should return the framework type of an attribute. This is named
	// differently than Type to prevent a conflict with the tfsdk.Attribute
	// field name.
	GetType() attr.Type

	// IsComputed should return true if the attribute configuration value is
	// computed. This is named differently than Computed to prevent a conflict
	// with the tfsdk.Attribute field name.
	IsComputed() bool

	// IsOptional should return true if the attribute configuration value is
	// optional. This is named differently than Optional to prevent a conflict
	// with the tfsdk.Attribute field name.
	IsOptional() bool

	// IsRequired should return true if the attribute configuration value is
	// required. This is named differently than Required to prevent a conflict
	// with the tfsdk.Attribute field name.
	IsRequired() bool

	// IsSensitive should return true if the attribute configuration value is
	// sensitive. This is named differently than Sensitive to prevent a
	// conflict with the tfsdk.Attribute field name.
	IsSensitive() bool

	// IsWriteOnly should return true if the attribute configuration value is
	// write-only. This is named differently than WriteOnly to prevent a
	// conflict with the tfsdk.Attribute field name.
	//
	// Write-only attributes are a managed-resource schema concept only.
	IsWriteOnly() bool
}
type SimplifiedSchema interface {
	TypeAtTerraformPath(context.Context, *tftypes.AttributePath) (attr.Type, error)
	// AttributeAtTerraformPath(context.Context, *tftypes.AttributePath) (AttributeSimplified, error)
}

func AttributePathValue(ctx context.Context, diags *diag.Diagnostics, attributePath *tftypes.AttributePath, src tfsdk.State, schema SimplifiedSchema) (attr.Value, path.Path) {
	convertedPath, localDiags := AttributePath(ctx, attributePath, schema)
	attrType, err := schema.TypeAtTerraformPath(ctx, attributePath)
	if err != nil {
		diags.AddError("Unable to get type for attribute path", fmt.Sprintf("%s: %s", attributePath.String(), err))
		return nil, convertedPath
	}
	if localDiags.HasError() {
		diags.Append(localDiags...)
		return nil, convertedPath
	}
	attrValue := attrType.ValueType(ctx)
	if localDiags := src.GetAttribute(ctx, convertedPath, &attrValue); localDiags.HasError() {
		diags.Append(localDiags...)
		return nil, convertedPath
	}
	return attrValue, convertedPath
}

func AttributePath(ctx context.Context, tfType *tftypes.AttributePath, schema SimplifiedSchema) (path.Path, diag.Diagnostics) {
	fwPath := path.Empty()

	for tfTypeStepIndex, tfTypeStep := range tfType.Steps() {
		currentTfTypeSteps := tfType.Steps()[:tfTypeStepIndex+1]
		currentTfTypePath := tftypes.NewAttributePathWithSteps(currentTfTypeSteps)
		attrType, err := schema.TypeAtTerraformPath(ctx, currentTfTypePath)

		if err != nil {
			return path.Empty(), diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						// Since this is an error with the attribute path
						// conversion, we cannot return a protocol path-based
						// diagnostic. Returning a framework human-readable
						// representation seems like the next best thing to do.
						fmt.Sprintf("Attribute Path: %s\n", currentTfTypePath.String())+
						fmt.Sprintf("Original Error: %s", err),
				),
			}
		}

		fwStep, err := AttributePathStep(ctx, tfTypeStep, attrType)

		if err != nil {
			return path.Empty(), diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is either an error in terraform-plugin-framework or a custom attribute type used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						// Since this is an error with the attribute path
						// conversion, we cannot return a protocol path-based
						// diagnostic. Returning a framework human-readable
						// representation seems like the next best thing to do.
						fmt.Sprintf("Attribute Path: %s\n", currentTfTypePath.String())+
						fmt.Sprintf("Original Error: %s", err),
				),
			}
		}

		// In lieu of creating a path.NewPathFromSteps function, this path
		// building logic is inlined to not expand the path package API.
		switch fwStep := fwStep.(type) {
		case path.PathStepAttributeName:
			fwPath = fwPath.AtName(string(fwStep))
		case path.PathStepElementKeyInt:
			fwPath = fwPath.AtListIndex(int(fwStep))
		case path.PathStepElementKeyString:
			fwPath = fwPath.AtMapKey(string(fwStep))
		case path.PathStepElementKeyValue:
			fwPath = fwPath.AtSetValue(fwStep.Value)
		default:
			return fwPath, diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is an error in terraform-plugin-framework used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						// Since this is an error with the attribute path
						// conversion, we cannot return a protocol path-based
						// diagnostic. Returning a framework human-readable
						// representation seems like the next best thing to do.
						fmt.Sprintf("Attribute Path: %s\n", currentTfTypePath.String())+
						fmt.Sprintf("Original Error: unknown path.PathStep type: %#v", fwStep),
				),
			}
		}
	}

	return fwPath, nil
}

func AttributePathStep(ctx context.Context, tfType tftypes.AttributePathStep, attrType attr.Type) (path.PathStep, error) {
	switch tfType := tfType.(type) {
	case tftypes.AttributeName:
		return path.PathStepAttributeName(string(tfType)), nil
	case tftypes.ElementKeyInt:
		return path.PathStepElementKeyInt(int64(tfType)), nil
	case tftypes.ElementKeyString:
		return path.PathStepElementKeyString(string(tfType)), nil
	case tftypes.ElementKeyValue:
		attrValue, err := Value(ctx, tftypes.Value(tfType), attrType)

		if err != nil {
			return nil, fmt.Errorf("unable to create PathStepElementKeyValue from tftypes.Value: %w", err)
		}

		return path.PathStepElementKeyValue{Value: attrValue}, nil
	default:
		return nil, fmt.Errorf("unknown tftypes.AttributePathStep: %#v", tfType)
	}
}

func Value(ctx context.Context, tfType tftypes.Value, attrType attr.Type) (attr.Value, error) {
	if attrType == nil {
		return nil, fmt.Errorf("unable to convert tftypes.Value (%s) to attr.Value: missing attr.Type", tfType.String())
	}

	attrValue, err := attrType.ValueFromTerraform(ctx, tfType)

	if err != nil {
		return nil, fmt.Errorf("unable to convert tftypes.Value (%s) to attr.Value: %w", tfType.String(), err)
	}

	return attrValue, nil
}
