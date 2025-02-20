package schemafunc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/stretchr/testify/assert"
)

func TestFindChanges(t *testing.T) {
	tests := map[string]struct {
		src      any
		dest     any
		expected []string
	}{
		"no changes": {
			src:      &TFSimpleModel{Name: types.StringValue("name")},
			dest:     &TFSimpleModel{Name: types.StringValue("name")},
			expected: []string{},
		},
		"simple change": {
			src:      &TFSimpleModel{Name: types.StringValue("name")},
			dest:     &TFSimpleModel{Name: types.StringValue("new-name")},
			expected: []string{"name"},
		},
		"object equal": {
			src: &TFSimpleModel{
				AdvancedConfig: asObjectValue(ctx, TFAdvancedConfig{JavascriptEnabled: types.BoolValue(true)}, AdvancedConfigObjType.AttrTypes),
			},
			dest: &TFSimpleModel{
				AdvancedConfig: asObjectValue(ctx, TFAdvancedConfig{JavascriptEnabled: types.BoolValue(true)}, AdvancedConfigObjType.AttrTypes),
			},
			expected: []string{},
		},
		"object change": {
			src: &TFSimpleModel{
				AdvancedConfig: asObjectValue(ctx, TFAdvancedConfig{JavascriptEnabled: types.BoolValue(true)}, AdvancedConfigObjType.AttrTypes),
			},
			dest: &TFSimpleModel{
				AdvancedConfig: asObjectValue(ctx, TFAdvancedConfig{JavascriptEnabled: types.BoolValue(false)}, AdvancedConfigObjType.AttrTypes),
			},
			expected: []string{"advanced_config", "advanced_config.javascript_enabled"},
		},
		"object change from null": {
			src: &TFSimpleModel{},
			dest: &TFSimpleModel{
				AdvancedConfig: asObjectValue(ctx, TFAdvancedConfig{JavascriptEnabled: types.BoolValue(false)}, AdvancedConfigObjType.AttrTypes),
			},
			expected: []string{"advanced_config", "advanced_config.javascript_enabled"},
		},
		"list equal": {
			src: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc}),
			},
			expected: []string{},
		},
		"list no change on unknown": {
			src: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigDest}),
			},
			expected: []string{},
		},
		"list change": {
			src: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone2"), []TFRegionConfig{regionConfigSrc}),
			},
			expected: []string{"replication_specs", "replication_specs[0]", "replication_specs[0].zone_name"},
		},
		"list add": {
			src: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc, regionConfigSrc}),
			},
			expected: []string{"replication_specs", "replication_specs[0]", "replication_specs[0].region_configs", "replication_specs[0].region_configs[1]"},
		},
		"list remove": {
			src: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc, regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc}),
			},
			expected: []string{"replication_specs", "replication_specs[0]", "replication_specs[0].region_configs", "replication_specs[0].region_configs[-1]"},
		},
		"list remove root": {
			src: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc, regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ReplicationSpecs: types.ListValueMust(ReplicationSpecsObjType, nil),
			},
			expected: []string{"replication_specs", "replication_specs[-0]"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := schemafunc.FindChanges(ctx, tc.src, tc.dest)
			assert.Equal(t, tc.expected, actual)
		})
	}
}