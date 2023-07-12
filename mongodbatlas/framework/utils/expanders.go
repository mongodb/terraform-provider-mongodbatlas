package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringSet accepts a `types.Set` and returns a slice of strings.
func StringSet(ctx context.Context, in types.Set) []string {
	results := []string{}
	_ = in.ElementsAs(ctx, &results, false)
	return results
}

// StringList accepts a `types.List` and returns a slice of strings.
func StringList(ctx context.Context, in types.List) []string {
	results := []string{}
	_ = in.ElementsAs(ctx, &results, false)
	return results
}
