package advancedcluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

// Object types mirroring the resource schema. Kept local to the test so building a model doesn't
// depend on unexported package state; ObjectValueFrom fails loudly if the schema and these diverge.
var (
	specsAttrTypes = map[string]attr.Type{
		"disk_iops":       types.Int64Type,
		"disk_size_gb":    types.Float64Type,
		"ebs_volume_type": types.StringType,
		"instance_size":   types.StringType,
		"node_count":      types.Int64Type,
	}
	autoScalingObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"compute_enabled":            types.BoolType,
		"compute_max_instance_size":  types.StringType,
		"compute_min_instance_size":  types.StringType,
		"compute_scale_down_enabled": types.BoolType,
		"disk_gb_enabled":            types.BoolType,
	}}
	specsObjectType       = types.ObjectType{AttrTypes: specsAttrTypes}
	regionConfigsObjeType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"analytics_auto_scaling": autoScalingObjectType,
		"analytics_specs":        specsObjectType,
		"auto_scaling":           autoScalingObjectType,
		"backing_provider_name":  types.StringType,
		"electable_specs":        specsObjectType,
		"priority":               types.Int64Type,
		"provider_name":          types.StringType,
		"read_only_specs":        specsObjectType,
		"region_name":            types.StringType,
	}}
	replicationSpecsObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"container_id":   types.MapType{ElemType: types.StringType},
		"external_id":    types.StringType,
		"region_configs": types.ListType{ElemType: regionConfigsObjeType},
		"zone_id":        types.StringType,
		"zone_name":      types.StringType,
	}}
)

// specs is a readable stand-in for TFSpecsModel in table tests; nil means the block is absent.
type specs struct {
	diskSizeGb    types.Float64
	ebsVolumeType types.String
	instanceSize  types.String
	diskIops      types.Int64
	nodeCount     types.Int64
}

func (s *specs) object(t *testing.T) types.Object {
	t.Helper()
	if s == nil {
		return types.ObjectNull(specsAttrTypes)
	}
	obj, diags := types.ObjectValueFrom(t.Context(), specsAttrTypes, advancedcluster.TFSpecsModel{
		DiskSizeGb:    nullIfZero(s.diskSizeGb, types.Float64Null()),
		DiskIops:      nullIfZero(s.diskIops, types.Int64Null()),
		EbsVolumeType: nullIfZero(s.ebsVolumeType, types.StringNull()),
		InstanceSize:  nullIfZero(s.instanceSize, types.StringNull()),
		NodeCount:     nullIfZero(s.nodeCount, types.Int64Null()),
	})
	require.False(t, diags.HasError(), "building specs object: %v", diags)
	return obj
}

// nullIfZero maps the struct zero value of an attr.Value to an explicit null of the right type.
func nullIfZero[T comparable](value, null T) T {
	var zero T
	if value == zero {
		return null
	}
	return value
}

// modelWithSpecs builds a TFModel holding a single replication spec with a single region config.
func modelWithSpecs(t *testing.T, analytics, electable, readOnly *specs) *advancedcluster.TFModel {
	t.Helper()
	ctx := t.Context()
	regionConfig, diags := types.ObjectValueFrom(ctx, regionConfigsObjeType.AttrTypes, advancedcluster.TFRegionConfigsModel{
		AnalyticsAutoScaling: types.ObjectNull(autoScalingObjectType.AttrTypes),
		AnalyticsSpecs:       analytics.object(t),
		AutoScaling:          types.ObjectNull(autoScalingObjectType.AttrTypes),
		BackingProviderName:  types.StringNull(),
		ElectableSpecs:       electable.object(t),
		Priority:             types.Int64Value(7),
		ProviderName:         types.StringValue("AZURE"),
		ReadOnlySpecs:        readOnly.object(t),
		RegionName:           types.StringValue("GERMANY_WEST_CENTRAL"),
	})
	require.False(t, diags.HasError(), "building region config: %v", diags)

	regionConfigs, diags := types.ListValue(regionConfigsObjeType, []attr.Value{regionConfig})
	require.False(t, diags.HasError(), "building region configs: %v", diags)

	repSpec, diags := types.ObjectValueFrom(ctx, replicationSpecsObjectType.AttrTypes, advancedcluster.TFReplicationSpecsModel{
		RegionConfigs: regionConfigs,
		ContainerId:   types.MapNull(types.StringType),
		ExternalId:    types.StringNull(),
		ZoneId:        types.StringNull(),
		ZoneName:      types.StringNull(),
	})
	require.False(t, diags.HasError(), "building replication spec: %v", diags)

	repSpecs, diags := types.ListValue(replicationSpecsObjectType, []attr.Value{repSpec})
	require.False(t, diags.HasError(), "building replication specs: %v", diags)

	return &advancedcluster.TFModel{
		MongoDBMajorVersion: types.StringNull(),
		Labels:              types.MapNull(types.StringType),
		Tags:                types.MapNull(types.StringType),
		ReplicationSpecs:    repSpecs,
	}
}

// analyticsSpecsOf reads back the analytics_specs of the first region config of the first rep spec.
func analyticsSpecsOf(t *testing.T, model *advancedcluster.TFModel) advancedcluster.TFSpecsModel {
	t.Helper()
	return specsOf(t, model, "analytics_specs")
}

func specsOf(t *testing.T, model *advancedcluster.TFModel, attrName string) advancedcluster.TFSpecsModel {
	t.Helper()
	ctx := t.Context()
	var repSpecs []advancedcluster.TFReplicationSpecsModel
	require.False(t, model.ReplicationSpecs.ElementsAs(ctx, &repSpecs, false).HasError())
	require.Len(t, repSpecs, 1)
	var regionConfigs []advancedcluster.TFRegionConfigsModel
	require.False(t, repSpecs[0].RegionConfigs.ElementsAs(ctx, &regionConfigs, false).HasError())
	require.Len(t, regionConfigs, 1)

	obj := map[string]types.Object{
		"analytics_specs": regionConfigs[0].AnalyticsSpecs,
		"electable_specs": regionConfigs[0].ElectableSpecs,
		"read_only_specs": regionConfigs[0].ReadOnlySpecs,
	}[attrName]

	var out advancedcluster.TFSpecsModel
	require.False(t, obj.As(ctx, &out, basetypes.ObjectAsOptions{}).HasError())
	return out
}

// analyticsSpecsIsNull reports whether Atlas' analytics_specs object stayed absent.
func analyticsSpecsIsNull(t *testing.T, model *advancedcluster.TFModel) bool {
	t.Helper()
	ctx := t.Context()
	var repSpecs []advancedcluster.TFReplicationSpecsModel
	require.False(t, model.ReplicationSpecs.ElementsAs(ctx, &repSpecs, false).HasError())
	require.Len(t, repSpecs, 1)
	var regionConfigs []advancedcluster.TFRegionConfigsModel
	require.False(t, repSpecs[0].RegionConfigs.ElementsAs(ctx, &regionConfigs, false).HasError())
	require.Len(t, regionConfigs, 1)
	return regionConfigs[0].AnalyticsSpecs.IsNull()
}

// TestAdvancedCluster_overrideUnreportedInactiveSpecs covers
// https://github.com/mongodb/terraform-provider-mongodbatlas/issues/4238: Atlas omits hardware
// attributes for a spec that describes no nodes, which used to be written to state as null and
// contradicted the configured value, surfacing as "Provider produced inconsistent result after apply".
func TestAdvancedCluster_overrideUnreportedInactiveSpecs(t *testing.T) {
	testCases := map[string]struct {
		prev           *specs // plan in Create/Update, prior state in Read
		atlas          *specs // what Atlas returned
		expected       *specs
		expectedIsNull bool
	}{
		"issue 4238: analytics disk_size_gb omitted for a spec with no nodes is kept from the config": {
			prev: &specs{
				diskSizeGb:   types.Float64Value(2048),
				diskIops:     types.Int64Value(7500),
				instanceSize: types.StringValue("M60"),
				nodeCount:    types.Int64Value(0),
			},
			atlas: &specs{
				instanceSize: types.StringValue("M60"),
				nodeCount:    types.Int64Value(0),
			},
			expected: &specs{
				diskSizeGb:   types.Float64Value(2048),
				diskIops:     types.Int64Value(7500),
				instanceSize: types.StringValue("M60"),
				nodeCount:    types.Int64Value(0),
			},
		},
		"every omitted attribute of an inactive spec is kept": {
			prev: &specs{
				diskSizeGb:    types.Float64Value(2048),
				diskIops:      types.Int64Value(7500),
				ebsVolumeType: types.StringValue("STANDARD"),
				instanceSize:  types.StringValue("M60"),
				nodeCount:     types.Int64Value(0),
			},
			atlas: &specs{nodeCount: types.Int64Value(0)},
			expected: &specs{
				diskSizeGb:    types.Float64Value(2048),
				diskIops:      types.Int64Value(7500),
				ebsVolumeType: types.StringValue("STANDARD"),
				instanceSize:  types.StringValue("M60"),
				nodeCount:     types.Int64Value(0),
			},
		},
		"Atlas values win over the previous model when it reports them": {
			prev: &specs{
				diskSizeGb:   types.Float64Value(2048),
				instanceSize: types.StringValue("M60"),
				nodeCount:    types.Int64Value(0),
			},
			atlas: &specs{
				diskSizeGb:   types.Float64Value(100),
				instanceSize: types.StringValue("M50"),
				nodeCount:    types.Int64Value(0),
			},
			expected: &specs{
				diskSizeGb:   types.Float64Value(100),
				instanceSize: types.StringValue("M50"),
				nodeCount:    types.Int64Value(0),
			},
		},
		"a spec with nodes is left untouched: Atlas reports the hardware actually backing them": {
			prev: &specs{
				diskSizeGb:   types.Float64Value(2048),
				instanceSize: types.StringValue("M60"),
				nodeCount:    types.Int64Value(3),
			},
			atlas: &specs{
				instanceSize: types.StringValue("M60"),
				nodeCount:    types.Int64Value(3),
			},
			expected: &specs{
				instanceSize: types.StringValue("M60"),
				nodeCount:    types.Int64Value(3),
			},
		},
		"unknown plan values never reach the state": {
			prev: &specs{
				diskSizeGb: types.Float64Unknown(),
				nodeCount:  types.Int64Value(0),
			},
			atlas:    &specs{nodeCount: types.Int64Value(0)},
			expected: &specs{nodeCount: types.Int64Value(0)},
		},
		"node_count omitted by Atlas falls back to the previous model": {
			prev: &specs{
				diskSizeGb: types.Float64Value(2048),
				nodeCount:  types.Int64Value(0),
			},
			atlas: &specs{},
			expected: &specs{
				diskSizeGb: types.Float64Value(2048),
				nodeCount:  types.Int64Value(0),
			},
		},
		"a spec absent from Atlas' response is not resurrected": {
			prev: &specs{
				diskSizeGb: types.Float64Value(2048),
				nodeCount:  types.Int64Value(0),
			},
			atlas:          nil,
			expectedIsNull: true,
		},
		"a spec absent from the previous model is left as Atlas reported it": {
			prev:     nil,
			atlas:    &specs{nodeCount: types.Int64Value(0)},
			expected: &specs{nodeCount: types.Int64Value(0)},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			modelIn := modelWithSpecs(t, tc.prev, nil, nil)
			modelOut := modelWithSpecs(t, tc.atlas, nil, nil)
			var diags diag.Diagnostics

			advancedcluster.OverrideAttributesWithPrevStateValue(t.Context(), modelIn, modelOut, &diags)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			if tc.expectedIsNull {
				assert.True(t, analyticsSpecsIsNull(t, modelOut), "analytics_specs should stay null")
				return
			}
			got := analyticsSpecsOf(t, modelOut)
			want := tc.expected
			assert.Equal(t, nullIfZero(want.diskSizeGb, types.Float64Null()), got.DiskSizeGb, "disk_size_gb")
			assert.Equal(t, nullIfZero(want.diskIops, types.Int64Null()), got.DiskIops, "disk_iops")
			assert.Equal(t, nullIfZero(want.ebsVolumeType, types.StringNull()), got.EbsVolumeType, "ebs_volume_type")
			assert.Equal(t, nullIfZero(want.instanceSize, types.StringNull()), got.InstanceSize, "instance_size")
			assert.Equal(t, nullIfZero(want.nodeCount, types.Int64Null()), got.NodeCount, "node_count")
		})
	}
}

// TestAdvancedCluster_overrideUnreportedInactiveSpecsAllSpecTypes checks the override applies to
// electable_specs and read_only_specs too, not only to the analytics_specs from the issue report.
func TestAdvancedCluster_overrideUnreportedInactiveSpecsAllSpecTypes(t *testing.T) {
	inactiveWithDisk := &specs{diskSizeGb: types.Float64Value(2048), nodeCount: types.Int64Value(0)}
	reportedWithoutDisk := &specs{nodeCount: types.Int64Value(0)}

	modelIn := modelWithSpecs(t, inactiveWithDisk, inactiveWithDisk, inactiveWithDisk)
	modelOut := modelWithSpecs(t, reportedWithoutDisk, reportedWithoutDisk, reportedWithoutDisk)
	var diags diag.Diagnostics

	advancedcluster.OverrideAttributesWithPrevStateValue(t.Context(), modelIn, modelOut, &diags)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	for _, attrName := range []string{"analytics_specs", "electable_specs", "read_only_specs"} {
		got := specsOf(t, modelOut, attrName)
		assert.Equal(t, types.Float64Value(2048), got.DiskSizeGb, "%s.disk_size_gb", attrName)
	}
}

// TestAdvancedCluster_overrideUnreportedInactiveSpecsNoReplicationSpecs guards the flex cluster path,
// where the model carries no replication specs to reconcile.
func TestAdvancedCluster_overrideUnreportedInactiveSpecsNoReplicationSpecs(t *testing.T) {
	for name, repSpecs := range map[string]types.List{
		"null":    types.ListNull(replicationSpecsObjectType),
		"unknown": types.ListUnknown(replicationSpecsObjectType),
		"empty":   types.ListValueMust(replicationSpecsObjectType, []attr.Value{}),
	} {
		t.Run(name, func(t *testing.T) {
			newModel := func() *advancedcluster.TFModel {
				return &advancedcluster.TFModel{
					MongoDBMajorVersion: types.StringNull(),
					Labels:              types.MapNull(types.StringType),
					Tags:                types.MapNull(types.StringType),
					ReplicationSpecs:    repSpecs,
				}
			}
			modelIn, modelOut := newModel(), newModel()
			var diags diag.Diagnostics

			advancedcluster.OverrideAttributesWithPrevStateValue(t.Context(), modelIn, modelOut, &diags)

			assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
			assert.Equal(t, repSpecs, modelOut.ReplicationSpecs)
		})
	}
}
