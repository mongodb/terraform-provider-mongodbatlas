package customplanmodifier_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
	"github.com/stretchr/testify/assert"
)

func TestAttributeChanges_AttributeChanged(t *testing.T) {
	tests := map[string]struct {
		attr     string
		changes  customplanmodifier.AttributeChanges
		expected bool
	}{
		"match found": {
			changes:  []string{"name", "description"},
			attr:     "name",
			expected: true,
		},
		"match not found": {
			changes:  []string{"name", "description"},
			attr:     "type",
			expected: false,
		},
		"nested attribute match": {
			changes:  []string{"config.name", "settings.enabled"},
			attr:     "name",
			expected: true,
		},
		"empty changes": {
			changes:  []string{},
			attr:     "name",
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.changes.AttributeChanged(tc.attr)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
func TestAttributeChanges_KeepUnknown(t *testing.T) {
	tests := map[string]struct {
		changes                  customplanmodifier.AttributeChanges
		attributeEffectedMapping map[string][]string
		expectedKeepUnknownAttrs []string
	}{
		"empty mapping": {
			changes:                  []string{"name", "description"},
			attributeEffectedMapping: map[string][]string{},
			expectedKeepUnknownAttrs: []string{},
		},
		"single mapping with match": {
			changes: []string{"name", "config.type"},
			attributeEffectedMapping: map[string][]string{
				"name": {"id", "status"},
			},
			expectedKeepUnknownAttrs: []string{"name", "id", "status"},
		},
		"multiple mappings with matches": {
			changes: []string{"name", "type", "config.value"},
			attributeEffectedMapping: map[string][]string{
				"name": {"id"},
				"type": {"category", "version"},
			},
			expectedKeepUnknownAttrs: []string{"name", "id", "type", "category", "version"},
		},
		"no matching changes": {
			changes: []string{"description", "status"},
			attributeEffectedMapping: map[string][]string{
				"name": {"id"},
				"type": {"category"},
			},
			expectedKeepUnknownAttrs: []string{},
		},
		"nested attribute changes": {
			changes: []string{"config.name", "settings.enabled"},
			attributeEffectedMapping: map[string][]string{
				"name":    {"id", "status"},
				"enabled": {"auth_status"},
			},
			expectedKeepUnknownAttrs: []string{"name", "id", "status", "enabled", "auth_status"},
		},
		"list attribute changes": {
			changes: []string{"replication_specs[0].zone_name"},
			attributeEffectedMapping: map[string][]string{
				"zone_name": {"priority", "region"},
			},
			expectedKeepUnknownAttrs: []string{"zone_name", "priority", "region"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.changes.KeepUnknown(tc.attributeEffectedMapping)
			assert.ElementsMatch(t, tc.expectedKeepUnknownAttrs, actual)
		})
	}
}

func TestAttributeChanges_ListIndexChanged(t *testing.T) {
	tests := map[string]struct {
		name     string
		changes  customplanmodifier.AttributeChanges
		index    int
		expected bool
	}{
		"empty changes": {
			name:     "replication_specs",
			index:    0,
			changes:  []string{},
			expected: false,
		},
		"list element modified": {
			name:     "replication_specs",
			index:    0,
			changes:  []string{"replication_specs[0]", "replication_specs[0].zone_name"},
			expected: true,
		},
		"list element added": {
			name:     "replication_specs",
			index:    0,
			changes:  []string{"replication_specs[+0]"},
			expected: false,
		},
		"list element removed": {
			name:     "replication_specs",
			index:    1,
			changes:  []string{"replication_specs[-1]"},
			expected: false,
		},
		"different index": {
			name:     "replication_specs",
			index:    1,
			changes:  []string{"replication_specs[0]", "replication_specs[0].zone_name"},
			expected: false,
		},
		"different list name": {
			name:     "other_specs",
			index:    0,
			changes:  []string{"replication_specs[0]", "replication_specs[0].zone_name"},
			expected: false,
		},
		"nested list": {
			name:     "replication_specs[0].region_configs",
			index:    0,
			changes:  []string{"replication_specs[0].region_configs[0]", "replication_specs[0].region_configs[0].priority"},
			expected: true,
		},
		"nested list false": {
			name:     "replication_specs[0].region_configs",
			index:    1,
			changes:  []string{"replication_specs[0].region_configs[0]", "replication_specs[0].region_configs[0].priority"},
			expected: false,
		},
		"index beyond bounds": {
			name:     "replication_specs",
			index:    5,
			changes:  []string{"replication_specs[0]", "replication_specs[1]"},
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.changes.ListIndexChanged(tc.name, tc.index)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
func TestAttributeChanges_NestedListLenChanges(t *testing.T) {
	tests := map[string]struct {
		fullPath string
		changes  customplanmodifier.AttributeChanges
		expected bool
	}{
		"empty changes": {
			fullPath: "replication_specs.region_configs",
			changes:  []string{},
			expected: false,
		},
		"no nested list changes": {
			fullPath: "replication_specs.region_configs",
			changes:  []string{"name", "description", "replication_specs.zone_name"},
			expected: false,
		},
		"add nested element": {
			fullPath: "replication_specs.region_configs",
			changes:  []string{"replication_specs.region_configs[+0]", "replication_specs.region_configs.priority"},
			expected: true,
		},
		"add nested element add different index should be false": {
			fullPath: "replication_specs[0].region_configs",
			changes:  []string{"replication_specs[1].region_configs[+0]"},
			expected: false,
		},
		"remove nested element": {
			fullPath: "replication_specs.region_configs",
			changes:  []string{"replication_specs.region_configs[-1]", "replication_specs.region_configs.region_name"},
			expected: true,
		},
		"mixed list operations": {
			fullPath: "replication_specs.region_configs",
			changes: []string{
				"replication_specs.region_configs[+0]",
				"replication_specs.region_configs[-1]",
				"replication_specs.region_configs.priority",
			},
			expected: true,
		},
		"different path": {
			fullPath: "other.configs",
			changes:  []string{"replication_specs.region_configs[+0]", "replication_specs.region_configs[-1]"},
			expected: false,
		},
		"multiple nested levels": {
			fullPath: "replication_specs.region_configs.zones",
			changes:  []string{"replication_specs.region_configs.zones[+0]", "replication_specs.region_configs[0].zones.name"},
			expected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.changes.ListLenChanges(tc.fullPath)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
