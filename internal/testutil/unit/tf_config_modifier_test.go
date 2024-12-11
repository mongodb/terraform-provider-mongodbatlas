package unit_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/require"
)

func TestApplyConfigModifiers(t *testing.T) {
	var nameReplacement = unit.TFConfigReplacement{
		Type:          unit.TFConfigReplacementString,
		ResourceName:  "cluster",
		AttributeName: "name",
	}
	tests := map[string]struct {
		oldConfig string
		newConfig string
		expected  string
		modifiers []unit.TFConfigReplacement
	}{
		"Replace attribute value": {
			oldConfig: `
resource "mongodbatlas_cluster" "test" {
  name = "old_name"
}`,
			newConfig: `
resource "mongodbatlas_cluster" "test" {
  name = "new_name"
  untouched = "yes"
}`,
			modifiers: []unit.TFConfigReplacement{nameReplacement},
			expected: `
resource "mongodbatlas_cluster" "test" {
  name = "old_name"
  untouched = "yes"
}`,
		},
		"No changes when oldConfig is empty": {
			oldConfig: ``,
			newConfig: `
resource "mongodbatlas_cluster" "test" {
  name = "new_name"
}`,
			modifiers: []unit.TFConfigReplacement{nameReplacement},
			expected:  ``,
		},
		"No changes when newConfig is empty": {
			oldConfig: `
resource "mongodbatlas_cluster" "test" {
  name = "old_name"
}`,
			newConfig: ``,
			modifiers: []unit.TFConfigReplacement{nameReplacement},
			expected:  ``,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := unit.ApplyConfigModifiers(t, tc.oldConfig, tc.newConfig, tc.modifiers)
			require.Equal(t, tc.expected, result)
		})
	}
}
