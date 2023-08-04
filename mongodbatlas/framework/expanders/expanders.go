package mongodbatlas

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
