package advancedclustertpf

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

func useStateForUnknown(ctx context.Context, diags *diag.Diagnostics, plan, state *TFModel, keepUnknown []string) {
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown)
	if slices.Contains(keepUnknown, replicationSpecsTFModelName) { // Early return if replication_specs is in keepUnknown
		return
	}
	// Nested fields are not supported by CopyUnknowns unless the whole field is Unknown.
	// Therefore, we need to handle replication_specs and partially unknown fields such as region_configs.(electable_specs|auto_scaling) manually.
	stateReplicationSpecs, planReplicationSpecs := statePlanTFListModel[TFReplicationSpecsModel](ctx, state.ReplicationSpecs, plan.ReplicationSpecs, diags)
	if diags.HasError() || len(planReplicationSpecs) != len(stateReplicationSpecs) {
		return
	}
	for i := range planReplicationSpecs {
		replicationSpecState := &stateReplicationSpecs[i]
		replicationSpecPlan := &planReplicationSpecs[i]
		schemafunc.CopyUnknowns(ctx, replicationSpecState, replicationSpecPlan, keepUnknown)
		fillInUnknownsInRegionConfigs(ctx, replicationSpecState, replicationSpecPlan, diags)
		if diags.HasError() {
			return
		}
	}
	newReplicationSpecs, localDiags := types.ListValueFrom(ctx, ReplicationSpecsObjType, planReplicationSpecs)
	diags.Append(localDiags...)
	if diags.HasError() {
		return
	}
	plan.ReplicationSpecs = newReplicationSpecs
}

func fillInUnknownsInRegionConfigs(ctx context.Context, replicationSpecState, replicationSpecPlan *TFReplicationSpecsModel, diags *diag.Diagnostics) {
	stateRegionConfigs, planRegionConfigs := statePlanTFListModel[TFRegionConfigsModel](ctx, replicationSpecState.RegionConfigs, replicationSpecPlan.RegionConfigs, diags)
	if diags.HasError() || len(stateRegionConfigs) != len(planRegionConfigs) {
		return
	}
	for j := range planRegionConfigs {
		stateRegionConfig := &stateRegionConfigs[j]
		planRegionConfig := &planRegionConfigs[j]
		useStateForUnknownRegionConfig(ctx, planRegionConfig, stateRegionConfig, diags)
		if diags.HasError() {
			return
		}
	}
	newRegionConfig, localDiags := types.ListValueFrom(ctx, RegionConfigsObjType, planRegionConfigs)
	if localDiags.HasError() {
		diags.Append(localDiags...)
		return
	}
	replicationSpecPlan.RegionConfigs = newRegionConfig
}

func useStateForUnknownRegionConfig(ctx context.Context, plan, state *TFRegionConfigsModel, diags *diag.Diagnostics) {
	if objectIsSet(plan.ElectableSpecs) {
		newElectableSpecs := copyUnknownsFromObject(ctx, diags, state.ElectableSpecs, plan.ElectableSpecs, func(ctx context.Context, newValue TFSpecsModel) (basetypes.ObjectValue, diag.Diagnostics) {
			return types.ObjectValueFrom(ctx, SpecsObjType.AttrTypes, newValue)
		})
		if diags.HasError() {
			return
		}
		plan.ElectableSpecs = newElectableSpecs
	}
	if objectIsSet(plan.AutoScaling) {
		newAutoScaling := copyUnknownsFromObject(ctx, diags, state.AutoScaling, plan.AutoScaling, func(ctx context.Context, newValue TFAutoScalingModel) (basetypes.ObjectValue, diag.Diagnostics) {
			return types.ObjectValueFrom(ctx, AutoScalingObjType.AttrTypes, newValue)
		})
		if diags.HasError() {
			return
		}
		plan.AutoScaling = newAutoScaling
	}
	schemafunc.CopyUnknowns(ctx, state, plan, []string{})
}

func copyUnknownsFromObject[T any](ctx context.Context, diags *diag.Diagnostics, state, plan types.Object, structToObject func(context.Context, T) (basetypes.ObjectValue, diag.Diagnostics)) basetypes.ObjectValue {
	planElectableSpecs := asTFObjectModel[T](ctx, plan, diags)
	stateElectableSpecs := asTFObjectModel[T](ctx, state, diags)
	if diags.HasError() {
		return basetypes.ObjectValue{}
	}
	schemafunc.CopyUnknowns(ctx, &stateElectableSpecs, &planElectableSpecs, []string{})
	newElectableSpecs, localDiags := structToObject(ctx, planElectableSpecs)
	if localDiags.HasError() {
		diags.Append(localDiags...)
		return basetypes.ObjectValue{}
	}
	return newElectableSpecs
}

func objectIsSet(obj types.Object) bool {
	return !obj.IsNull() && !obj.IsUnknown()
}

func statePlanTFListModel[T any](ctx context.Context, stateList, planList types.List, diags *diag.Diagnostics) (state, plan []T) {
	stateElements := asTFListModel[T](ctx, stateList, diags)
	planElements := asTFListModel[T](ctx, planList, diags)
	return stateElements, planElements
}

func asTFListModel[T any](ctx context.Context, list types.List, diags *diag.Diagnostics) []T {
	elements := make([]T, len(list.Elements()))
	if localDiags := list.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return elements
}

func asTFObjectModel[T any](ctx context.Context, obj types.Object, diags *diag.Diagnostics) T {
	var element T
	if localDiags := obj.As(ctx, &element, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return element
	}
	return element
}
