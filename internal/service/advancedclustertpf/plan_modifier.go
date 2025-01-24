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

func IsUnknown(obj reflect.Value) *bool {
	method := obj.MethodByName("IsUnknown")
	if !method.IsValid() { // Method not found
		return nil
	}
	results := method.Call([]reflect.Value{})
	if len(results) > 0 {
		result := results[0]
		response, ok := result.Interface().(bool)
		if ok {
			return &response
		}
	}
	return nil
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
		isUnknownP := IsUnknown(field)
		if isUnknownP != nil && *isUnknownP {
			return true
		}
	}
	return false
}

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
		if !found {
			continue
		}
		isUnknownP := IsUnknown(valDest.Field(i))
		if isUnknownP == nil || !*isUnknownP {
			continue
		}
		if !valDest.Field(i).CanSet() {
			continue
		}
		tflog.Info(ctx, fmt.Sprintf("Copying unknown field: %s", name))
		valDest.Field(i).Set(valSrc.FieldByName(name))
	}
}

func useRemoteForUnknown(ctx context.Context, diags *diag.Diagnostics, plan, remoteModel *TFModel, keepUnknown []string) {
	CopyUnknowns(ctx, remoteModel, plan, keepUnknown)
	if slices.Contains(keepUnknown, replicationSpecsTFModelName) {
		return
	}
	// Nested fields are not supported by CopyUnknowns unless the whole field is Unknown.
	// Therefore, we need to handle replication_specs and partially unknown fields such as region_configs.(electable_specs|auto_scaling) manually.
	planReplicationSpecsElements := plan.ReplicationSpecs.Elements()
	readModelReplicationSpecsElements := remoteModel.ReplicationSpecs.Elements()
	if len(planReplicationSpecsElements) != len(readModelReplicationSpecsElements) {
		return
	}
	planReplicationSpecs := make([]TFReplicationSpecsModel, len(planReplicationSpecsElements))
	if localDiags := plan.ReplicationSpecs.ElementsAs(ctx, &planReplicationSpecs, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return
	}
	remoteReplicationSpecs := make([]TFReplicationSpecsModel, len(readModelReplicationSpecsElements))
	if localDiags := remoteModel.ReplicationSpecs.ElementsAs(ctx, &remoteReplicationSpecs, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return
	}
	for i := range planReplicationSpecs {
		replicationSpecRemote := &remoteReplicationSpecs[i]
		replicationSpecPlan := &planReplicationSpecs[i]
		CopyUnknowns(ctx, replicationSpecRemote, replicationSpecPlan, keepUnknown)
		fillInUnknownsInRegionConfigs(ctx, replicationSpecRemote, replicationSpecPlan, diags)
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

func fillInUnknownsInRegionConfigs(ctx context.Context, replicationSpecRemote, replicationSpecPlan *TFReplicationSpecsModel, diags *diag.Diagnostics) {
	regionConfigsRemoteElements := replicationSpecRemote.RegionConfigs.Elements()
	regionConfigsPlanElements := replicationSpecPlan.RegionConfigs.Elements()
	if len(regionConfigsRemoteElements) != len(regionConfigsPlanElements) {
		return
	}
	remoteRegionConfigs := make([]TFRegionConfigsModel, len(regionConfigsRemoteElements))
	if localDiags := replicationSpecRemote.RegionConfigs.ElementsAs(ctx, &remoteRegionConfigs, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return
	}
	planRegionConfigs := make([]TFRegionConfigsModel, len(regionConfigsPlanElements))
	if localDiags := replicationSpecPlan.RegionConfigs.ElementsAs(ctx, &planRegionConfigs, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return
	}
	for j := range regionConfigsPlanElements {
		remoteRegionConfig := &remoteRegionConfigs[j]
		planRegionConfig := &planRegionConfigs[j]
		if !planRegionConfig.ElectableSpecs.IsNull() && !planRegionConfig.ElectableSpecs.IsUnknown() {
			planElectableSpecs := &TFSpecsModel{}
			if localDiags := planRegionConfig.ElectableSpecs.As(ctx, planElectableSpecs, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
				diags.Append(localDiags...)
				return
			}
			remoteElectableSpecs := &TFSpecsModel{}
			if localDiags := remoteRegionConfig.ElectableSpecs.As(ctx, remoteElectableSpecs, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
				diags.Append(localDiags...)
				return
			}
			CopyUnknowns(ctx, remoteElectableSpecs, planElectableSpecs, []string{})
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
			autoScalingRemote := &TFAutoScalingModel{}
			if localDiags := remoteRegionConfig.AutoScaling.As(ctx, autoScalingRemote, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
				diags.Append(localDiags...)
				return
			}
			CopyUnknowns(ctx, autoScalingRemote, autoScalingPlan, []string{})
			newAutoScaling, localDiags := types.ObjectValueFrom(ctx, AutoScalingObjType.AttrTypes, autoScalingPlan)
			if localDiags.HasError() {
				diags.Append(localDiags...)
				return
			}
			planRegionConfig.AutoScaling = newAutoScaling
		}
		CopyUnknowns(ctx, remoteRegionConfig, planRegionConfig, []string{})
	}
	newRegionConfig, localDiags := types.ListValueFrom(ctx, RegionConfigsObjType, planRegionConfigs)
	if localDiags.HasError() {
		diags.Append(localDiags...)
		return
	}
	replicationSpecPlan.RegionConfigs = newRegionConfig
}
