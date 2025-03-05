package advancedclustertpf

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

var (
	attributeRootChangeMapping = map[string][]string{
		"disk_size_gb":           {}, // disk_size_gb can be change at any level/spec
		"replication_specs":      {},
		"mongo_db_major_version": {"mongo_db_version"},
	}
	attributeReplicationSpecChangeMapping = map[string][]string{
		// All these fields can exist in specs that are computed, therefore, it is not safe to use them when they have changed.
		"disk_iops":       {},
		"ebs_volume_type": {},
		"disk_size_gb":    {},                  // disk_size_gb can be change at any level/spec
		"instance_size":   {"disk_iops"},       // disk_iops can change based on instance_size changes
		"provider_name":   {"ebs_volume_type"}, // AWS --> AZURE will change ebs_volume_type
		"region_name":     {"container_id"},    // container_id changes based on region_name changes
		"zone_name":       {"zone_id"},         // zone_id copy from state is not safe when
	}
	autoScalingBoolValues   = []string{"compute_enabled", "disk_gb_enabled", "compute_scale_down_enabled"}
	autoScalingStringValues = []string{"compute_min_instance_size", "compute_max_instance_size"}
	keepUnknownsCalls       = []func(string, attr.Value) bool{
		// when node_count != 0 --> keepUnknown
		func(name string, replacement attr.Value) bool {
			return name == "node_count" && !replacement.Equal(types.Int64Value(0))
		},
		// when auto_scaling bool attributes are true --> keepUnknown
		func(name string, replacement attr.Value) bool {
			return slices.Contains(autoScalingBoolValues, name) && replacement.Equal(types.BoolValue(true))
		},
		// when auto_scaling string attributes are non empty, (M10/M30) --> keepUnknown
		func(name string, replacement attr.Value) bool {
			return slices.Contains(autoScalingStringValues, name) && replacement.(types.String).ValueString() != ""
		},
	}
)

// useStateForUnknowns should be called only in Update, because of findClusterDiff
func useStateForUnknowns(ctx context.Context, diags *diag.Diagnostics, cfg, state, plan *TFModel) {
	diff := findClusterDiff(ctx, state, cfg, plan, diags)
	if diags.HasError() || diff.isAnyUpgrade() { // Don't do anything in upgrades
		return
	}
	attributeChanges := schemafunc.NewAttributeChanges(ctx, state, plan)
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	keepUnknown = append(keepUnknown, attributeChanges.KeepUnknown(attributeRootChangeMapping)...)
	keepUnknown = append(keepUnknown, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown, keepUnknownsCalls...)
	if slices.Contains(keepUnknown, "replication_specs") {
		useStateForUnknownsReplicationSpecs(ctx, diags, state, plan, &attributeChanges)
	}
}

func useStateForUnknownsReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, attrChanges *schemafunc.AttributeChanges) {
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return
	}
	planWithUnknowns := []TFReplicationSpecsModel{}
	keepUnknownsUnchangedSpec := determineKeepUnknownsUnchangedReplicationSpecs(ctx, diags, state, plan, attrChanges)
	keepUnknownsUnchangedSpec = append(keepUnknownsUnchangedSpec, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
	if diags.HasError() {
		return
	}
	for i := range planRepSpecsTF {
		if i < len(stateRepSpecsTF) {
			keepUnknowns := keepUnknownsUnchangedSpec
			if attrChanges.ListIndexChanged("replication_specs", i) {
				keepUnknowns = determineKeepUnknownsChangedReplicationSpec(keepUnknownsUnchangedSpec, attrChanges, fmt.Sprintf("replication_specs[%d]", i))
			}
			schemafunc.CopyUnknowns(ctx, &stateRepSpecsTF[i], &planRepSpecsTF[i], keepUnknowns, keepUnknownsCalls...)
		}
		planWithUnknowns = append(planWithUnknowns, planRepSpecsTF[i])
	}
	listType, diagsLocal := types.ListValueFrom(ctx, ReplicationSpecsObjType, planWithUnknowns)
	diags.Append(diagsLocal...)
	if diags.HasError() {
		return
	}
	plan.ReplicationSpecs = listType
}

// determineKeepUnknownsChangedReplicationSpec: These fields must be kept unknown in the replication_specs[index_of_changes]
func determineKeepUnknownsChangedReplicationSpec(keepUnknownsAlways []string, attributeChanges *schemafunc.AttributeChanges, parentPath string) []string {
	var keepUnknowns = slices.Clone(keepUnknownsAlways)
	if attributeChanges.NestedListLenChanges(parentPath + ".region_configs") {
		keepUnknowns = append(keepUnknowns, "container_id")
	}
	return append(keepUnknowns, attributeChanges.KeepUnknown(attributeReplicationSpecChangeMapping)...)
}

func determineKeepUnknownsUnchangedReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, attributeChanges *schemafunc.AttributeChanges) []string {
	keepUnknowns := []string{}
	// Could be set to "" if we are using an ISS cluster
	if usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags) { // When using new sharding config, the legacy id must never be copied
		keepUnknowns = append(keepUnknowns, "id")
	}
	// for isShardingConfigUpgrade, it will be empty in the plan, so we need to keep it unknown
	// for listLenChanges, it might be an insertion in the middle of replication spec leading to wrong value from state copied
	if isShardingConfigUpgrade(ctx, state, plan, diags) || attributeChanges.ListLenChanges("replication_specs") {
		keepUnknowns = append(keepUnknowns, "external_id")
	}
	return keepUnknowns
}

func determineKeepUnknownsAutoScaling(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) []string {
	var keepUnknown []string
	computedUsed, diskUsed := autoScalingUsed(ctx, diags, state, plan)
	if computedUsed {
		keepUnknown = append(keepUnknown, "instance_size")
		keepUnknown = append(keepUnknown, attributeReplicationSpecChangeMapping["instance_size"]...)
	}
	if diskUsed {
		keepUnknown = append(keepUnknown, "disk_size_gb")
		keepUnknown = append(keepUnknown, attributeReplicationSpecChangeMapping["disk_size_gb"]...)
	}
	return keepUnknown
}

// autoScalingUsed checks is auto-scaling was enabled (state) or will be enabled (plan).
func autoScalingUsed(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) (computedUsed, diskUsed bool) {
	for _, model := range []*TFModel{state, plan} {
		repSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, model.ReplicationSpecs)
		for i := range repSpecsTF {
			regiongConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, repSpecsTF[i].RegionConfigs)
			for j := range regiongConfigsTF {
				for _, autoScalingTF := range []types.Object{regiongConfigsTF[j].AutoScaling, regiongConfigsTF[j].AnalyticsAutoScaling} {
					if autoScalingTF.IsNull() || autoScalingTF.IsUnknown() {
						continue
					}
					autoscaling := TFModelObject[TFAutoScalingModel](ctx, diags, autoScalingTF)
					if autoscaling == nil {
						continue
					}
					if autoscaling.ComputeEnabled.ValueBool() {
						computedUsed = true
					}
					if autoscaling.DiskGBEnabled.ValueBool() {
						diskUsed = true
					}
				}
			}
		}
	}
	return
}

// setExplicitEmptyAutoScaling sets the auto-scaling to a known value with false for enabled flags
func setExplicitEmptyAutoScaling(ctx context.Context, diags *diag.Diagnostics, plan *TFModel, planResp *tfsdk.Plan) {
	repSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	for i := range repSpecsTF {
		regiongConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, repSpecsTF[i].RegionConfigs)
		for j := range regiongConfigsTF {
			autoScalingEmptyValues := asObjectValue(ctx, TFAutoScalingModel{
				ComputeEnabled:          types.BoolValue(false),
				ComputeMinInstanceSize:  types.StringValue(""),
				ComputeMaxInstanceSize:  types.StringValue(""),
				ComputeScaleDownEnabled: types.BoolValue(false),
				DiskGBEnabled:           types.BoolValue(false),
			}, AutoScalingObjType.AttrTypes)
			autoScalingPath := path.Root("replication_specs").AtListIndex(i).AtName("region_configs").AtListIndex(j).AtName("auto_scaling")
			tflog.Info(ctx, fmt.Sprintf("Setting auto-scaling to empty values, for path %s", autoScalingPath))
			localDiags := planResp.SetAttribute(ctx, autoScalingPath, autoScalingEmptyValues)
			if localDiags.HasError() {
				tflog.Error(ctx, fmt.Sprintf("Failed to set auto-scaling to empty values: %v", localDiags))
			}
			diags.Append(localDiags...)
		}
	}
}

func TFModelList[T any](ctx context.Context, diags *diag.Diagnostics, input types.List) []T {
	elements := make([]T, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return elements
}
func asObjectValue[T any](ctx context.Context, t T, attrs map[string]attr.Type) types.Object {
	objType, diagsLocal := types.ObjectValueFrom(ctx, attrs, t)
	if diagsLocal.HasError() {
		panic("failed to convert object to model")
	}
	return objType
}

func TFModelObject[T any](ctx context.Context, diags *diag.Diagnostics, input types.Object) *T {
	item := new(T)
	if localDiags := input.As(ctx, item, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return item
}

func triggerConfigChanges(ctx context.Context, diags *diag.Diagnostics, config, state, plan *TFModel, planResp *tfsdk.Plan) {
	triggerWhenAutoScalingRemoved(ctx, diags, state, plan, config, planResp)
	triggerWhenSpecBlockRemoved(ctx, diags, state, config, planResp)
}

func triggerWhenAutoScalingRemoved(ctx context.Context, diags *diag.Diagnostics, state, plan, config *TFModel, planResp *tfsdk.Plan) {
	computeUsed, diskUsed := autoScalingUsed(ctx, diags, state, plan)
	if !computeUsed && !diskUsed {
		return
	}
	configComputeUsed, configDiskUsed := autoScalingUsed(ctx, diags, config, config)
	if configComputeUsed || configDiskUsed {
		return
	}
	if computeUsed || diskUsed {
		setExplicitEmptyAutoScaling(ctx, diags, plan, planResp)
	}
}

func triggerWhenSpecBlockRemoved(ctx context.Context, diags *diag.Diagnostics, state, cfg *TFModel, planResp *tfsdk.Plan) {
	repSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	repSpecsConfigTF := TFModelList[TFReplicationSpecsModel](ctx, diags, cfg.ReplicationSpecs)
	for i := range repSpecsTF {
		if i >= len(repSpecsConfigTF) {
			continue
		}
		regiongConfigsStateTF := TFModelList[TFRegionConfigsModel](ctx, diags, repSpecsTF[i].RegionConfigs)
		regiongConfigsCfgTF := TFModelList[TFRegionConfigsModel](ctx, diags, repSpecsConfigTF[i].RegionConfigs)
		for j := range regiongConfigsStateTF {
			if j >= len(regiongConfigsCfgTF) {
				continue
			}
			specsState := map[string]*TFSpecsModel{
				"analytics_specs": TFModelObject[TFSpecsModel](ctx, diags, regiongConfigsStateTF[j].AnalyticsSpecs),
				"electable_specs": TFModelObject[TFSpecsModel](ctx, diags, regiongConfigsStateTF[j].ElectableSpecs),
				"read_only_specs": TFModelObject[TFSpecsModel](ctx, diags, regiongConfigsStateTF[j].ReadOnlySpecs),
			}
			specsCfg := map[string]types.Object{
				"analytics_specs": regiongConfigsCfgTF[j].AnalyticsSpecs,
				"electable_specs": regiongConfigsCfgTF[j].ElectableSpecs,
				"read_only_specs": regiongConfigsCfgTF[j].ReadOnlySpecs,
			}
			for name, specState := range specsState {
				if specState == nil || specState.NodeCount.ValueInt64() == 0 || !specsCfg[name].IsNull() {
					continue
				}
				specPath := path.Root("replication_specs").AtListIndex(i).AtName("region_configs").AtListIndex(j).AtName(name)
				tflog.Info(ctx, fmt.Sprintf("Setting %s to empty values, for path %s", name, specPath))
				emptySpec := asObjectValue(ctx, TFSpecsModel{
					NodeCount:     types.Int64Value(0),
					DiskSizeGb:    types.Float64Unknown(),
					EbsVolumeType: types.StringUnknown(),
					InstanceSize:  types.StringUnknown(),
					DiskIops:      types.Int64Unknown(),
				}, SpecsObjType.AttrTypes)
				localDiags := planResp.SetAttribute(ctx, specPath, emptySpec)
				if localDiags.HasError() {
					tflog.Error(ctx, fmt.Sprintf("Failed to set %s to empty values: %v", name, localDiags))
				}
				diags.Append(localDiags...)
			}
		}
	}
}
