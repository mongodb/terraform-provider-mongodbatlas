package autogen

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type QueryParamArg struct {
	Value   any
	APIName string
}

// BuildQueryParamMap builds a query parameter map from Terraform values.
// It supports types.String, types.Int64, types.Bool, types.List, and types.Set.
// For List and Set types, element types of String, Int64, and Bool are supported.
// Values that are null or unknown are skipped.
func BuildQueryParamMap(ctx context.Context, args []QueryParamArg) map[string]string {
	result := make(map[string]string)

	for _, arg := range args {
		var value string
		var isSet bool

		switch v := arg.Value.(type) {
		case types.String:
			if !v.IsNull() && !v.IsUnknown() {
				value = v.ValueString()
				isSet = true
			}
		case types.Int64:
			if !v.IsNull() && !v.IsUnknown() {
				value = fmt.Sprintf("%d", v.ValueInt64())
				isSet = true
			}
		case types.Bool:
			if !v.IsNull() && !v.IsUnknown() {
				value = fmt.Sprintf("%t", v.ValueBool())
				isSet = true
			}
		case types.List:
			if !v.IsNull() && !v.IsUnknown() {
				value, isSet = extractListValue(ctx, v)
			}
		case types.Set:
			if !v.IsNull() && !v.IsUnknown() {
				value, isSet = extractSetValue(ctx, v)
			}
		}

		if isSet {
			result[arg.APIName] = value
		}
	}

	return result
}

// extractListValue extracts values from a types.List and joins them with commas.
// Supports element types: String, Int64, and Bool.
func extractListValue(ctx context.Context, list types.List) (string, bool) {
	if len(list.Elements()) == 0 {
		return "", false
	}

	elementType := list.ElementType(ctx)
	values := extractCollectionValues(ctx, list.Elements(), elementType, "list")

	if len(values) == 0 {
		return "", false
	}
	return strings.Join(values, ","), true
}

// extractSetValue extracts values from a types.Set and joins them with commas.
// Supports element types: String, Int64, and Bool.
func extractSetValue(ctx context.Context, set types.Set) (string, bool) {
	if len(set.Elements()) == 0 {
		return "", false
	}

	elementType := set.ElementType(ctx)
	values := extractCollectionValues(ctx, set.Elements(), elementType, "set")

	if len(values) == 0 {
		return "", false
	}
	return strings.Join(values, ","), true
}

// extractCollectionValues extracts values from a collection of elements based on their type.
// Supports String, Int64, and Bool element types.
func extractCollectionValues(ctx context.Context, elements []attr.Value, elementType attr.Type, collectionType string) []string {
	var values []string

	for _, elem := range elements {
		var strValue string
		var isValid bool

		switch elementType {
		case types.StringType:
			if v, ok := elem.(types.String); ok && !v.IsNull() && !v.IsUnknown() {
				strValue = v.ValueString()
				isValid = true
			}
		case types.Int64Type:
			if v, ok := elem.(types.Int64); ok && !v.IsNull() && !v.IsUnknown() {
				strValue = fmt.Sprintf("%d", v.ValueInt64())
				isValid = true
			}
		case types.BoolType:
			if v, ok := elem.(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
				strValue = fmt.Sprintf("%t", v.ValueBool())
				isValid = true
			}
		default:
			tflog.Warn(ctx, "Unsupported element type in query parameter collection",
				map[string]any{
					"collection_type": collectionType,
					"element_type":    fmt.Sprintf("%T", elementType),
				})
			continue
		}

		if isValid {
			values = append(values, strValue)
		}
	}

	return values
}
