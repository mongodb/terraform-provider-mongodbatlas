package advancedclustertpf

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func IsUnknown(obj reflect.Value) bool {
	method := obj.MethodByName("IsUnknown")
	if !method.IsValid() {
		panic(fmt.Sprintf("IsUnknown method not found for %v", obj))
	}
	results := method.Call([]reflect.Value{})
	if len(results) != 1 {
		panic(fmt.Sprintf("IsUnknown method must return a single value, got %v", results))
	}
	result := results[0]
	response, ok := result.Interface().(bool)
	if !ok {
		panic(fmt.Sprintf("IsUnknown method must return a bool, got %v", result))
	}
	return response
}

func HasUnknowns(obj any) bool {
	valObj := reflect.ValueOf(obj)
	if valObj.Kind() != reflect.Ptr {
		panic("params must be pointer")
	}
	valObj = valObj.Elem()
	if valObj.Kind() != reflect.Struct {
		panic("params must be pointer to struct")
	}
	typeObj := valObj.Type()
	for i := range typeObj.NumField() {
		field := valObj.Field(i)
		if IsUnknown(field) {
			return true
		}
	}
	return false
}

// CopyUnknowns use reflection to copy unknown fields from src to dest.
// The alternative without reflection would need to pass every field in a struct.
// The implementation is similar to internal/common/conversion/model_generation.go#CopyModel
func CopyUnknowns(ctx context.Context, src, dest any, keepUnknown []string) {
	valSrc := reflect.ValueOf(src)
	valDest := reflect.ValueOf(dest)
	if valSrc.Kind() != reflect.Ptr || valDest.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("params must be pointers %v %v\n", src, dest))
	}
	valSrc = valSrc.Elem()
	valDest = valDest.Elem()
	if valSrc.Kind() != reflect.Struct || valDest.Kind() != reflect.Struct {
		panic(fmt.Sprintf("params must be pointers to structs: %v %v\n", src, dest))
	}
	typeSrc := valSrc.Type()
	typeDest := valDest.Type()
	for i := range typeDest.NumField() {
		fieldDest := typeDest.Field(i)
		name := fieldDest.Name
		if slices.Contains(keepUnknown, name) {
			continue
		}
		_, found := typeSrc.FieldByName(name)
		if !found || !IsUnknown(valDest.Field(i)) || !valDest.Field(i).CanSet() {
			continue
		}
		tflog.Info(ctx, fmt.Sprintf("Copying unknown field: %s", name))
		valDest.Field(i).Set(valSrc.FieldByName(name))
	}
}

func useStateForUnknown(ctx context.Context, diags *diag.Diagnostics, plan, state *TFModel, keepUnknown []string) {
	CopyUnknowns(ctx, state, plan, keepUnknown)
	if slices.Contains(keepUnknown, replicationSpecsTFModelName) { // Early return if replication_specs is in keepUnknown
		return
	}
	// Nested fields are not supported by CopyUnknowns unless the whole field is Unknown.
	// Therefore, we need to handle replication_specs and partially unknown fields such as region_configs.(electable_specs|auto_scaling) manually.
	planReplicationSpecsElements := plan.ReplicationSpecs.Elements()
	readModelReplicationSpecsElements := state.ReplicationSpecs.Elements()
	if len(planReplicationSpecsElements) != len(readModelReplicationSpecsElements) {
		return
	}
	planReplicationSpecs := make([]TFReplicationSpecsModel, len(planReplicationSpecsElements))
	if localDiags := plan.ReplicationSpecs.ElementsAs(ctx, &planReplicationSpecs, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return
	}
	stateReplicationSpecs := make([]TFReplicationSpecsModel, len(readModelReplicationSpecsElements))
	if localDiags := state.ReplicationSpecs.ElementsAs(ctx, &stateReplicationSpecs, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return
	}
	for i := range planReplicationSpecs {
		replicationSpecState := &stateReplicationSpecs[i]
		replicationSpecPlan := &planReplicationSpecs[i]
		CopyUnknowns(ctx, replicationSpecState, replicationSpecPlan, keepUnknown)
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
	regionConfigsStateElements := replicationSpecState.RegionConfigs.Elements()
	regionConfigsPlanElements := replicationSpecPlan.RegionConfigs.Elements()
	if len(regionConfigsStateElements) != len(regionConfigsPlanElements) {
		return
	}
	stateRegionConfigs := make([]TFRegionConfigsModel, len(regionConfigsStateElements))
	if localDiags := replicationSpecState.RegionConfigs.ElementsAs(ctx, &stateRegionConfigs, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return
	}
	planRegionConfigs := make([]TFRegionConfigsModel, len(regionConfigsPlanElements))
	if localDiags := replicationSpecPlan.RegionConfigs.ElementsAs(ctx, &planRegionConfigs, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return
	}
	for j := range regionConfigsPlanElements {
		stateRegionConfig := &stateRegionConfigs[j]
		planRegionConfig := &planRegionConfigs[j]
		if !planRegionConfig.ElectableSpecs.IsNull() && !planRegionConfig.ElectableSpecs.IsUnknown() {
			planElectableSpecs := &TFSpecsModel{}
			if localDiags := planRegionConfig.ElectableSpecs.As(ctx, planElectableSpecs, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
				diags.Append(localDiags...)
				return
			}
			stateElectableSpecs := &TFSpecsModel{}
			if localDiags := stateRegionConfig.ElectableSpecs.As(ctx, stateElectableSpecs, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
				diags.Append(localDiags...)
				return
			}
			CopyUnknowns(ctx, stateElectableSpecs, planElectableSpecs, []string{})
			newElectableSpecs, localDiags := types.ObjectValueFrom(ctx, SpecsObjType.AttrTypes, planElectableSpecs)
			if localDiags.HasError() {
				diags.Append(localDiags...)
				return
			}
			planRegionConfig.ElectableSpecs = newElectableSpecs
		}
		if !planRegionConfig.AutoScaling.IsNull() && !planRegionConfig.AutoScaling.IsUnknown() {
			autoScalingPlan := &TFAutoScalingModel{}
			if localDiags := planRegionConfig.AutoScaling.As(ctx, autoScalingPlan, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
				diags.Append(localDiags...)
				return
			}
			autoScalingState := &TFAutoScalingModel{}
			if localDiags := stateRegionConfig.AutoScaling.As(ctx, autoScalingState, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
				diags.Append(localDiags...)
				return
			}
			CopyUnknowns(ctx, autoScalingState, autoScalingPlan, []string{})
			newAutoScaling, localDiags := types.ObjectValueFrom(ctx, AutoScalingObjType.AttrTypes, autoScalingPlan)
			if localDiags.HasError() {
				diags.Append(localDiags...)
				return
			}
			planRegionConfig.AutoScaling = newAutoScaling
		}
		CopyUnknowns(ctx, stateRegionConfig, planRegionConfig, []string{})
	}
	newRegionConfig, localDiags := types.ListValueFrom(ctx, RegionConfigsObjType, planRegionConfigs)
	if localDiags.HasError() {
		diags.Append(localDiags...)
		return
	}
	replicationSpecPlan.RegionConfigs = newRegionConfig
}
