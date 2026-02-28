package conversion

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ToAnySlicePointer converts to a slice pointer of any as needed in some Atlas SDK Go structs
func ToAnySlicePointer(value *[]map[string]any) *[]any {
	if value == nil {
		return nil
	}
	ret := make([]any, len(*value))
	for i, item := range *value {
		ret[i] = item
	}
	return &ret
}

func TFSetValueOrNull[T any](ctx context.Context, ptr *[]T, elemType attr.Type) types.Set {
	if ptr == nil || len(*ptr) == 0 {
		return types.SetNull(elemType)
	}
	set, _ := types.SetValueFrom(ctx, elemType, *ptr)
	return set
}
