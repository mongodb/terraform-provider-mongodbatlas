package advancedclustertpf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
)

func manualPlanChanges(ctx context.Context, diags *diag.Diagnostics, differ *customplanmodifier.PlanModifyDiffer) (manualChanges bool) {
	if nonZeroSpecRemoved(ctx, diags, differ) {
		manualChanges = true
	}
	if didRemoveOrChangeAutoScaling(ctx, diags, differ, "auto_scaling") {
		manualChanges = true
	}
	if didRemoveOrChangeAutoScaling(ctx, diags, differ, "analytics_auto_scaling") {
		manualChanges = true
	}
	return manualChanges
}

func nonZeroSpecRemoved(ctx context.Context, diags *diag.Diagnostics, differ *customplanmodifier.PlanModifyDiffer) (manualChanges bool) {
	removedSpecs := func(diffs []customplanmodifier.DiffTPF[TFSpecsModel]) []customplanmodifier.DiffTPF[TFSpecsModel] {
		var removed []customplanmodifier.DiffTPF[TFSpecsModel]
		for _, diff := range diffs {
			if !diff.Removed() {
				continue
			}
			stateValue := diff.State
			nodeCount := stateValue.NodeCount
			if nodeCount.IsNull() || nodeCount.Equal(types.Int64Value(0)) {
				continue
			}
			removed = append(removed, diff)
		}
		return removed
	}
	analyticsSpecs := customplanmodifier.StateConfigDiffs[TFSpecsModel](ctx, diags, differ, "analytics_specs", false)
	readOnlySpecs := customplanmodifier.StateConfigDiffs[TFSpecsModel](ctx, diags, differ, "read_only_specs", false)
	if diags.HasError() {
		return false
	}
	for _, spec := range removedSpecs(analyticsSpecs) {
		manualChanges = true
		tflog.Info(ctx, fmt.Sprintf("AnalyticsSpecs @ %s removed\n%s", spec.Path.String(), spec.State))
		explicitRemoveSpec := TFSpecsModel{
			NodeCount:     types.Int64Value(0),
			InstanceSize:  spec.State.InstanceSize,
			DiskIops:      types.Int64Unknown(),
			EbsVolumeType: types.StringUnknown(),
			DiskSizeGb:    types.Float64Unknown(),
		}
		customplanmodifier.UpdatePlanValue(ctx, diags, differ, spec.Path, asObjectValue(ctx, explicitRemoveSpec, SpecsObjType.AttrTypes))
	}
	for _, spec := range removedSpecs(readOnlySpecs) {
		manualChanges = true
		tflog.Info(ctx, fmt.Sprintf("ReadOnlySpecs @ %s removed\n%s", spec.Path.String(), spec.State))
		electableSpecPath := spec.Path.ParentPath().AtName("electable_specs")
		electableSpec := customplanmodifier.ReadPlanStructValue[TFSpecsModel](ctx, diags, differ, electableSpecPath)
		if diags.HasError() {
			return false
		}
		explicitRemoveSpec := TFSpecsModel{
			NodeCount:     types.Int64Value(0),
			InstanceSize:  electableSpec.InstanceSize, // Use electable spec instance size in case it is updated
			DiskIops:      types.Int64Unknown(),
			EbsVolumeType: types.StringUnknown(),
			DiskSizeGb:    types.Float64Unknown(),
		}
		customplanmodifier.UpdatePlanValue(ctx, diags, differ, spec.Path, asObjectValue(ctx, explicitRemoveSpec, SpecsObjType.AttrTypes))
	}
	return manualChanges
}

var boolFalse = types.BoolValue(false)

func didRemoveOrChangeAutoScaling(ctx context.Context, diags *diag.Diagnostics, differ *customplanmodifier.PlanModifyDiffer, name string) (removedFlag bool) {
	autoScalings := customplanmodifier.StateConfigDiffs[TFAutoScalingModel](ctx, diags, differ, name, true)
	if diags.HasError() {
		return false
	}
	for _, autoScaling := range autoScalings {
		var explicitRemoveAutoScaling *TFAutoScalingModel
		if autoScaling.Removed() {
			if autoScaling.State.ComputeEnabled.Equal(boolFalse) && autoScaling.State.DiskGBEnabled.Equal(boolFalse) {
				continue
			}
			tflog.Info(ctx, fmt.Sprintf("AutoScaling @ %s removed\n%s", autoScaling.Path.String(), autoScaling.State))
			if name == "auto_scaling" {
				explicitRemoveAutoScaling = &TFAutoScalingModel{
					ComputeEnabled:          types.BoolValue(false),
					DiskGBEnabled:           types.BoolValue(false),
					ComputeMinInstanceSize:  types.StringUnknown(),
					ComputeMaxInstanceSize:  types.StringUnknown(),
					ComputeScaleDownEnabled: types.BoolValue(false),
				}
			} else { // analytics_auto_scaling is null from the backend by default
				explicitRemoveAutoScaling = nil
			}
		} else {
			explicitRemoveAutoScaling = autoScalingAttributeRemoved(ctx, autoScaling)
			if explicitRemoveAutoScaling == nil {
				continue
			}
		}
		removedFlag = true
		customplanmodifier.UpdatePlanValue(ctx, diags, differ, autoScaling.Path, asObjectValue(ctx, explicitRemoveAutoScaling, AutoScalingObjType.AttrTypes))
	}
	return removedFlag
}

func autoScalingAttributeRemoved(ctx context.Context, autoScaling customplanmodifier.DiffTPF[TFAutoScalingModel]) *TFAutoScalingModel {
	if !autoScaling.Changed() {
		return nil
	}
	stateValue := autoScaling.State
	configValue := autoScaling.Config
	var attributesRemoved []string
	if stateValue.ComputeEnabled.Equal(types.BoolValue(true)) && configValue.ComputeEnabled.IsNull() {
		attributesRemoved = append(attributesRemoved, "compute_enabled")
		configValue.ComputeEnabled = types.BoolValue(false)
	}
	if stateValue.DiskGBEnabled.Equal(types.BoolValue(true)) && configValue.DiskGBEnabled.IsNull() {
		attributesRemoved = append(attributesRemoved, "disk_gb_enabled")
		configValue.DiskGBEnabled = types.BoolValue(false)
	}
	if stateValue.ComputeScaleDownEnabled.Equal(types.BoolValue(true)) && configValue.ComputeScaleDownEnabled.IsNull() {
		attributesRemoved = append(attributesRemoved, "compute_scale_down_enabled")
		configValue.ComputeScaleDownEnabled = types.BoolValue(false)
	}
	if len(attributesRemoved) == 0 {
		return nil
	}
	tflog.Info(ctx, fmt.Sprintf("Removed attributes %s of auto_scaling @ %s", strings.Join(attributesRemoved, ","), autoScaling.Path.String()))
	return configValue
}
