package conversion

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20241023001/admin"
)

func TestNewResourceTags(t *testing.T) {
	testCases := map[string]struct {
		plan     types.Map
		expected *[]admin.ResourceTag
	}{
		"tags null":    {types.MapNull(types.StringType), &[]admin.ResourceTag{}},
		"tags unknown": {types.MapUnknown(types.StringType), &[]admin.ResourceTag{}},
		"tags convert normally": {types.MapValueMust(types.StringType, map[string]attr.Value{
			"key1": types.StringValue("value1"),
		}), &[]admin.ResourceTag{
			*admin.NewResourceTag("key1", "value1"),
		}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, NewResourceTags(context.Background(), tc.plan))
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
			assert.Equal(t, tc.expected, NewTFTags(tc.adminTags))
		})
	}
}
