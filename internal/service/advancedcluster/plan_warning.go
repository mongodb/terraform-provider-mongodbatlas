package advancedcluster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func warnIgnoredSpecChanges(ctx context.Context, diags *diag.Diagnostics, config tfsdk.Config, state, plan *TFModel) {
	if !plan.UseEffectiveFields.ValueBool() {
		return
	}
	hardwareFields := []string{"instance_size", "disk_size_gb", "disk_iops"}
	var ignoredPaths []string
	stateRegions := singleReplicationSpecRegions(ctx, state.ReplicationSpecs)
	for i, value := range singleReplicationSpecRegions(ctx, plan.ReplicationSpecs) {
		region := TFModelObject[TFRegionConfigsModel](ctx, value.(types.Object))
		if region == nil || !isKnown(region.ProviderName) || !isKnown(region.RegionName) {
			continue
		}
		for _, stateValue := range stateRegions {
			prior := TFModelObject[TFRegionConfigsModel](ctx, stateValue.(types.Object))
			// Match by location so reordered or newly added regions do not produce false warnings.
			if prior == nil || !prior.ProviderName.Equal(region.ProviderName) || !prior.RegionName.Equal(region.RegionName) {
				continue
			}
			regionPath := path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(i)
			compute, disk := unchangedAutoScaling(prior.AutoScaling, region.AutoScaling)
			if compute || disk {
				ignoredPaths = append(ignoredPaths, changedSpecPaths(ctx, diags, config, regionPath.AtName("electable_specs"), prior.ElectableSpecs, region.ElectableSpecs, hardwareFields...)...)
				ignoredPaths = append(ignoredPaths, changedSpecPaths(ctx, diags, config, regionPath.AtName("read_only_specs"), prior.ReadOnlySpecs, region.ReadOnlySpecs, hardwareFields...)...)
			}
			analyticsCompute, _ := unchangedAutoScaling(prior.AnalyticsAutoScaling, region.AnalyticsAutoScaling)
			if analyticsCompute {
				ignoredPaths = append(ignoredPaths, changedSpecPaths(ctx, diags, config, regionPath.AtName("analytics_specs"), prior.AnalyticsSpecs, region.AnalyticsSpecs, "instance_size")...)
			}
			break
		}
	}
	if diags.HasError() || len(ignoredPaths) == 0 {
		return
	}
	diags.AddAttributeWarning(path.Root("replication_specs"),
		"Spec changes are ignored while auto-scaling remains enabled",
		fmt.Sprintf("With use_effective_fields = true and auto-scaling remaining enabled, Atlas ignores changes to the following attributes, although Terraform stores their new values in state:\n\n- %s\n\n"+
			"To apply these changes, disable auto-scaling and apply the desired values, then re-enable auto-scaling in a separate apply. "+
			"See: https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#manually-updating-specs-with-use_effective_fields",
			strings.Join(ignoredPaths, "\n- ")))
}

func singleReplicationSpecRegions(ctx context.Context, specs types.List) []attr.Value {
	if len(specs.Elements()) != 1 {
		return nil
	}
	spec := TFModelObject[TFReplicationSpecsModel](ctx, specs.Elements()[0].(types.Object))
	if spec == nil {
		return nil
	}
	return spec.RegionConfigs.Elements()
}

// unchangedAutoScaling returns enabled flags only when neither flag changes; Atlas applies requested specs when either is toggled.
func unchangedAutoScaling(state, plan types.Object) (computeEnabled, diskEnabled bool) {
	if !isKnown(state) || !isKnown(plan) {
		return false, false
	}
	before, after := state.Attributes(), plan.Attributes()
	for _, field := range []string{"compute_enabled", "disk_gb_enabled"} {
		if before[field].IsUnknown() || after[field].IsUnknown() || !before[field].Equal(after[field]) {
			return false, false
		}
	}
	return after["compute_enabled"].(types.Bool).ValueBool(), after["disk_gb_enabled"].(types.Bool).ValueBool()
}

func changedSpecPaths(ctx context.Context, diags *diag.Diagnostics, config tfsdk.Config, specsPath path.Path, state, plan types.Object, fields ...string) []string {
	if !isKnown(state) || !isKnown(plan) {
		return nil
	}
	var configured types.Object
	diags.Append(config.GetAttribute(ctx, specsPath, &configured)...)
	if diags.HasError() || !isKnown(configured) {
		return nil
	}
	before, after, configuredAttrs := state.Attributes(), plan.Attributes(), configured.Attributes()
	var ignored []string
	for _, field := range fields {
		if !isKnown(configuredAttrs[field]) || !isKnown(before[field]) || !isKnown(after[field]) || before[field].Equal(after[field]) {
			continue
		}
		if field == "instance_size" && after[field].(types.String).ValueString() == "AUTO" {
			continue
		}
		ignored = append(ignored, specsPath.AtName(field).String())
	}
	return ignored
}
