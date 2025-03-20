package conversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/types@v1.13.0
func AsUnknownValue(ctx context.Context, value attr.Value) attr.Value {
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
