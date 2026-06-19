package advancedcluster

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// MINIMAL-CONFIG PROTOTYPE: minimal unit tests. Not exhaustive — they show the
// pure helpers, that a minimal INFINITE config synthesizes a full spec, and that an
// explicit config is left untouched (reverse-compat). End-to-end is covered by the
// `terraform plan` of the example configs.

func TestParseProviderRegion(t *testing.T) {
	cases := map[string]struct{ in, wantProvider, wantRegion string }{
		"AWS:US_EAST_1":               {"AWS:US_EAST_1", "AWS", "US_EAST_1"},
		"GCP:CENTRAL_US":              {"GCP:CENTRAL_US", "GCP", "CENTRAL_US"},
		"empty -> full default":       {"", "AWS", "US_EAST_1"},
		"no colon (malformed) -> def": {"GCP", "AWS", "US_EAST_1"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gotProvider, gotRegion := ParseProviderRegion(tc.in)
			assert.Equal(t, tc.wantProvider, gotProvider)
			assert.Equal(t, tc.wantRegion, gotRegion)
		})
	}
}

func TestDefaultInstanceSizeForProfile(t *testing.T) {
	assert.Equal(t, "M30", DefaultInstanceSizeForProfile(ClusterProfileInfinite))
	assert.Equal(t, "M10", DefaultInstanceSizeForProfile(ClusterProfileCore))
	assert.Equal(t, "M10", DefaultInstanceSizeForProfile(""), "unset profile behaves as CORE")
}

func TestApplyMinimalConfigDefaults_MinimalInfinite(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	const projectID = "60ddf0123456789012345678" // project_id is Required: user always supplies it
	// config: the now-optional inputs (cluster_type, replication_specs) omitted; project_id provided.
	config := &TFModel{
		ClusterType:      types.StringNull(),
		ProjectID:        types.StringValue(projectID),
		ReplicationSpecs: types.ListNull(replicationSpecsObjType),
	}
	// plan: user supplied project_id + profile + provider_region; omitted inputs are Unknown (Computed).
	plan := &TFModel{
		ClusterProfile:   types.StringValue(ClusterProfileInfinite),
		ProviderRegion:   types.StringValue("AWS:US_EAST_1"),
		ClusterType:      types.StringUnknown(),
		ProjectID:        types.StringValue(projectID),
		ReplicationSpecs: types.ListUnknown(replicationSpecsObjType),
	}

	changed := applyMinimalConfigDefaults(ctx, &diags, config, plan)
	assert.False(t, diags.HasError(), "diags: %v", diags.Errors())
	assert.True(t, changed)
	assert.Equal(t, "REPLICASET", plan.ClusterType.ValueString())
	assert.Equal(t, projectID, plan.ProjectID.ValueString(), "project_id is Required and left untouched")
	assert.False(t, plan.ReplicationSpecs.IsNull() || plan.ReplicationSpecs.IsUnknown())

	// The synthesized region config carries provider/region + INFINITE auto-scaling defaults.
	repSpecs := TFModelList[TFReplicationSpecsModel](ctx, &diags, plan.ReplicationSpecs)
	assert.Len(t, repSpecs, 1)
	rcs := TFModelList[TFRegionConfigsModel](ctx, &diags, repSpecs[0].RegionConfigs)
	assert.Len(t, rcs, 1)
	assert.Equal(t, "AWS", rcs[0].ProviderName.ValueString())
	assert.Equal(t, "US_EAST_1", rcs[0].RegionName.ValueString())

	es := TFModelObject[TFSpecsModel](ctx, rcs[0].ElectableSpecs)
	assert.Equal(t, "M30", es.InstanceSize.ValueString())
	assert.Equal(t, int64(3), es.NodeCount.ValueInt64())

	as := TFModelObject[TFAutoScalingModel](ctx, rcs[0].AutoScaling)
	assert.True(t, as.ComputeEnabled.ValueBool())
	assert.Equal(t, "M30", as.ComputeMinInstanceSize.ValueString())
	assert.Equal(t, "M50", as.ComputeMaxInstanceSize.ValueString())
}

func TestApplyMinimalConfigDefaults_ExplicitConfigUnchanged(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	// Reverse-compat: user set every (now-optional) field explicitly -> nothing changes.
	explicitSpecs := types.ListValueMust(replicationSpecsObjType, []attr.Value{}) // non-null
	config := &TFModel{
		ClusterType:      types.StringValue("SHARDED"),
		ProjectID:        types.StringValue("abcabcabcabcabcabcabcabc"),
		ReplicationSpecs: explicitSpecs,
	}
	plan := &TFModel{
		ClusterProfile:   types.StringValue(ClusterProfileInfinite), // even INFINITE must not override explicit input
		ClusterType:      types.StringValue("SHARDED"),
		ProjectID:        types.StringValue("abcabcabcabcabcabcabcabc"),
		ReplicationSpecs: explicitSpecs,
	}

	changed := applyMinimalConfigDefaults(ctx, &diags, config, plan)
	assert.False(t, diags.HasError())
	assert.False(t, changed, "explicit config must not be modified")
	assert.Equal(t, "SHARDED", plan.ClusterType.ValueString())
	assert.Equal(t, "abcabcabcabcabcabcabcabc", plan.ProjectID.ValueString())
}
