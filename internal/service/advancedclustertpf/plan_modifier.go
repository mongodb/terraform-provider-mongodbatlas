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
	planReplicationSpecs := asTFListModel[TFReplicationSpecsModel](ctx, plan.ReplicationSpecs, diags)
	stateReplicationSpecs := asTFListModel[TFReplicationSpecsModel](ctx, state.ReplicationSpecs, diags)
	if diags.HasError() || len(planReplicationSpecs) != len(stateReplicationSpecs) {
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
	stateRegionConfigs := asTFListModel[TFRegionConfigsModel](ctx, replicationSpecState.RegionConfigs, diags)
	planRegionConfigs := asTFListModel[TFRegionConfigsModel](ctx, replicationSpecPlan.RegionConfigs, diags)
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
		planElectableSpecs := asTFObjectModel[TFSpecsModel](ctx, plan.ElectableSpecs, diags)
		stateElectableSpecs := asTFObjectModel[TFSpecsModel](ctx, state.ElectableSpecs, diags)
		if diags.HasError() {
			return
		}
		CopyUnknowns(ctx, &stateElectableSpecs, &planElectableSpecs, []string{})
		newElectableSpecs, localDiags := types.ObjectValueFrom(ctx, SpecsObjType.AttrTypes, planElectableSpecs)
		if localDiags.HasError() {
			diags.Append(localDiags...)
			return
		}
		plan.ElectableSpecs = newElectableSpecs
	}
	if objectIsSet(plan.AutoScaling) {
		autoScalingPlan := asTFObjectModel[TFAutoScalingModel](ctx, plan.AutoScaling, diags)
		autoScalingState := asTFObjectModel[TFAutoScalingModel](ctx, state.AutoScaling, diags)
		if diags.HasError() {
			return
		}
		CopyUnknowns(ctx, &autoScalingState, &autoScalingPlan, []string{})
		newAutoScaling, localDiags := types.ObjectValueFrom(ctx, AutoScalingObjType.AttrTypes, autoScalingPlan)
		if localDiags.HasError() {
			diags.Append(localDiags...)
			return
		}
		plan.AutoScaling = newAutoScaling
	}
	CopyUnknowns(ctx, state, plan, []string{})
	return
}

func objectIsSet(obj types.Object) bool {
	return !obj.IsNull() && !obj.IsUnknown()
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