package conversion

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TypesSetToString(ctx context.Context, set types.Set) []string {
	results := []string{}
	_ = set.ElementsAs(ctx, &results, false)
	return results
}

func TypesListToString(ctx context.Context, set types.List) []string {
	results := []string{}
	_ = set.ElementsAs(ctx, &results, false)
	return results
}

// StringValueToFramework converts a string value to a Framework String value.
// An empty string is converted to a null String. Useful for optional attributes.
func StringValueToFramework(v string) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}
