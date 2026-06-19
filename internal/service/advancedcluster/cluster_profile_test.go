package advancedcluster_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

// cluster_profile PROTOTYPE: minimal unit tests demonstrating the two behaviors.
// Not exhaustive — just enough to show CORE = unchanged and INFINITE = defaults applied.

func TestInstanceSizeTwoTiersUp(t *testing.T) {
	testCases := map[string]struct {
		instanceSize string
		expected     string
	}{
		"M30 -> M50 (the documented example)": {instanceSize: "M30", expected: "M50"},
		"M10 -> M30":                          {instanceSize: "M10", expected: "M30"},
		"M60 -> M140 (non-uniform jump)":      {instanceSize: "M60", expected: "M140"},
		"M400 caps at top (M700)":             {instanceSize: "M400", expected: "M700"},
		"M700 already at top stays M700":      {instanceSize: "M700", expected: "M700"},
		"unknown size returned unchanged":     {instanceSize: "M0", expected: "M0"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, advancedcluster.InstanceSizeTwoTiersUp(tc.instanceSize))
		})
	}
}

func TestClusterProfileAutoScaling(t *testing.T) {
	t.Run("INFINITE without user auto_scaling applies the defaults", func(t *testing.T) {
		model, apply := advancedcluster.ClusterProfileAutoScaling(advancedcluster.ClusterProfileInfinite, "M30", false)
		assert.True(t, apply)
		assert.True(t, model.ComputeEnabled.ValueBool())
		assert.True(t, model.ComputeScaleDownEnabled.ValueBool())
		assert.Equal(t, "M30", model.ComputeMinInstanceSize.ValueString(), "min = configured instance size")
		assert.Equal(t, "M50", model.ComputeMaxInstanceSize.ValueString(), "max = two tiers up")
		assert.False(t, model.DiskGBEnabled.ValueBool(), "only compute auto-scaling is enabled")
	})

	t.Run("CORE is baseline: no defaults applied", func(t *testing.T) {
		_, apply := advancedcluster.ClusterProfileAutoScaling(advancedcluster.ClusterProfileCore, "M30", false)
		assert.False(t, apply)
	})

	t.Run("unset profile is baseline: no defaults applied", func(t *testing.T) {
		_, apply := advancedcluster.ClusterProfileAutoScaling("", "M30", false)
		assert.False(t, apply)
	})

	t.Run("INFINITE honors explicit user auto_scaling (explicit input wins)", func(t *testing.T) {
		_, apply := advancedcluster.ClusterProfileAutoScaling(advancedcluster.ClusterProfileInfinite, "M30", true)
		assert.False(t, apply)
	})

	t.Run("INFINITE with no configured instance size is a no-op", func(t *testing.T) {
		_, apply := advancedcluster.ClusterProfileAutoScaling(advancedcluster.ClusterProfileInfinite, "", false)
		assert.False(t, apply)
	})
}
