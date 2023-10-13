package conversion

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/util"
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

// StringNullIfEmpty converts a string value to a Framework String value.
// An empty string is converted to a null String. Useful for optional attributes.
func StringNullIfEmpty(v string) types.String {
	return StringPtrNullIfEmpty(&v)
}

// StringPtrNullIfEmpty is similar to StringNullIfEmpty but can also handle nil string pointers.
func StringPtrNullIfEmpty(p *string) types.String {
	if util.IsStringPresent(p) {
		return types.StringValue(*p)
	}
	return types.StringNull()
}
