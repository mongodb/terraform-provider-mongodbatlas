package conversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

func TestNewResourceTags(t *testing.T) {
	testCases := map[string]struct {
		expected *[]admin.ResourceTag
		plan     types.Map
	}{
		"tags null":    {&[]admin.ResourceTag{}, types.MapNull(types.StringType)},
		"tags unknown": {&[]admin.ResourceTag{}, types.MapUnknown(types.StringType)},
		"tags convert normally": {&[]admin.ResourceTag{
			*admin.NewResourceTag("key1", "value1"),
		}, types.MapValueMust(types.StringType, map[string]attr.Value{
			"key1": types.StringValue("value1"),
		})},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, conversion.NewResourceTags(t.Context(), tc.plan))
		})
	}
}

func TestNewTFTags(t *testing.T) {
	var (
		tfMapEmpty     = types.MapValueMust(types.StringType, map[string]attr.Value{})
		apiListEmpty   = []admin.ResourceTag{}
		apiSingleTag   = []admin.ResourceTag{*admin.NewResourceTag("key1", "value1")}
		tfMapSingleTag = types.MapValueMust(types.StringType, map[string]attr.Value{"key1": types.StringValue("value1")})
	)
	testCases := map[string]struct {
		expected  types.Map
		adminTags []admin.ResourceTag
	}{
		"api empty list tf null should give map null":      {tfMapEmpty, apiListEmpty},
		"tags single value tf null should give map single": {tfMapSingleTag, apiSingleTag},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, conversion.NewTFTags(tc.adminTags))
		})
	}
}
