package conversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func TestToAnySlicePointer(t *testing.T) {
	testCases := map[string]*[]map[string]any{
		"nil":         nil,
		"empty":       {},
		"one element": {{"hi": "there"}},
		"more complex": {
			{"hi": "there"},
			{"bye": 1234},
		},
	}
	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := conversion.ToAnySlicePointer(value)
			if ret == nil {
				assert.Nil(t, value)
			} else {
				assert.NotNil(t, ret)
				assert.Len(t, *value, len(*ret))
				for i := range *value {
					assert.Equal(t, (*value)[i], (*ret)[i])
				}
			}
		})
	}
}

func TestTFSetValueOrNull(t *testing.T) {
	ctx := t.Context()

	testCases := map[string]*[]string{
		"nil":       nil,
		"empty":     {},
		"populated": {"a", "b", "c"},
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			result := conversion.TFSetValueOrNull(ctx, value, types.StringType)
			if value == nil || len(*value) == 0 {
				assert.True(t, result.IsNull())
			} else {
				assert.False(t, result.IsNull())
				assert.Len(t, result.Elements(), len(*value))
			}
		})
	}
}
