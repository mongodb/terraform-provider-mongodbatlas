package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type SimplifiedSchema interface {
	TypeAtTerraformPath(context.Context, *tftypes.AttributePath) (attr.Type, error)
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

// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/types@v1.13.0
func asUnknownValue(ctx context.Context, value attr.Value) attr.Value {
	switch v := value.(type) {
	case types.List:
		return types.ListUnknown(v.ElementType(ctx))
	case types.Object:
		return types.ObjectUnknown(v.AttributeTypes(ctx))
	case types.Map:
		return types.MapUnknown(v.ElementType(ctx))
	case types.Set:
		return types.SetUnknown(v.ElementType(ctx))
	case types.Tuple:
		return types.TupleUnknown(v.ElementTypes(ctx))
	case types.String:
		return types.StringUnknown()
	case types.Bool:
		return types.BoolUnknown()
	case types.Int64:
		return types.Int64Unknown()
	case types.Int32:
		return types.Int32Unknown()
	case types.Float64:
		return types.Float64Unknown()
	case types.Float32:
		return types.Float32Unknown()
	case types.Number:
		return types.NumberUnknown()
	case types.Dynamic:
		return types.DynamicUnknown()
	}
	panic(fmt.Sprintf("Unknown value to create unknown: %v", value))
}
