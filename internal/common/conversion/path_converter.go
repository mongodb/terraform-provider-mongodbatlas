package conversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type TPFSchema interface {
	TypeAtTerraformPath(context.Context, *tftypes.AttributePath) (attr.Type, error)
}

type TPFSrc interface {
	GetAttribute(context.Context, path.Path, any) diag.Diagnostics
}

// AttributePathValue retrieves the value for src (state/plan/config) @ attributePath with converted path.Path, schema is needed to get the correct types.XXX (String/Object/etc.)
func AttributePathValue(ctx context.Context, diags *diag.Diagnostics, attributePath *tftypes.AttributePath, src TPFSrc, schema TPFSchema) (attr.Value, path.Path) {
	convertedPath, localDiags := AttributePath(ctx, attributePath, schema)
	if localDiags.HasError() {
		diags.Append(localDiags...)
		return nil, convertedPath
	}
	attrType, err := schema.TypeAtTerraformPath(ctx, attributePath)
	if err != nil {
		diags.AddError("Unable to get type for attribute path", fmt.Sprintf("%s: %s", attributePath.String(), err))
		return nil, convertedPath
	}
	attrValue := attrType.ValueType(ctx)
	if localDiags := src.GetAttribute(ctx, convertedPath, &attrValue); localDiags.HasError() {
		diags.Append(localDiags...)
		return nil, convertedPath
	}
	return attrValue, convertedPath
}

// AttributePath similar to the internal function in TPF, but simpler interface as argument and less logging
func AttributePath(ctx context.Context, tfType *tftypes.AttributePath, schema TPFSchema) (path.Path, diag.Diagnostics) {
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
