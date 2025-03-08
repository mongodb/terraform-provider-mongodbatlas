package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type DiffHelper struct {
	req             *resource.ModifyPlanRequest
	resp            *resource.ModifyPlanResponse
	stateConfigDiff []tftypes.ValueDiff
}

func ReadConfigValue[T any](ctx context.Context, diags *diag.Diagnostics, d *DiffHelper, p path.Path) *T {
	var obj types.Object
	if localDiags := d.req.Config.GetAttribute(ctx, p, &obj); localDiags.HasError() {
		diags.Append(localDiags...)
		return nil
	}
	return TFModelObject[T](ctx, diags, obj)
}

func UpdatePlanValue[T attr.Value](ctx context.Context, diags *diag.Diagnostics, d *DiffHelper, p path.Path, value T) {
	if localDiags := d.resp.Plan.SetAttribute(ctx, p, value); localDiags.HasError() {
		diags.Append(localDiags...)
	}
}

type DiffTPF[T any] struct {
	Path          path.Path
	Plan          *T
	State         *T
	Config        *T
	PlanUnknown   bool
	ConfigUnknown bool
}

func (d *DiffTPF[T]) Removed() bool {
	return d.State != nil && d.Config == nil
}

func (d *DiffTPF[T]) PlanOrStateValue() *T {
	if d.Plan != nil {
		return d.Plan
	}
	return d.State
}

func StateConfigDiffs[T any](ctx context.Context, diags *diag.Diagnostics, d *DiffHelper, name tftypes.AttributeName, schema SimplifiedSchema) []DiffTPF[T] {
	earlyReturn := func(localDiags diag.Diagnostics) []DiffTPF[T] {
		diags.Append(localDiags...)
		return nil
	}
	var diffs []DiffTPF[T]

	for _, diff := range d.stateConfigDiff {
		if diff.Path.LastStep().Equal(name) {
			p, localDiags := AttributePath(ctx, diff.Path, schema)
			if localDiags.HasError() {
				return earlyReturn(localDiags)
			}
			var stateObj, configObj, planObj types.Object
			if d1 := d.req.State.GetAttribute(ctx, p, &stateObj); d1.HasError() {
				return earlyReturn(d1)
			}
			if d2 := d.req.Config.GetAttribute(ctx, p, &configObj); d2.HasError() {
				return earlyReturn(d2)
			}
			if d3 := d.req.Plan.GetAttribute(ctx, p, &planObj); d3.HasError() {
				return earlyReturn(d3)
			}
			var configParsed, stateParsed, planParsed *T
			if !stateObj.IsNull() { // stateObj is never unknown
				stateParsed = TFModelObject[T](ctx, diags, stateObj)
			}
			if !configObj.IsNull() && !configObj.IsUnknown() {
				configParsed = TFModelObject[T](ctx, diags, configObj)
			}
			if !planObj.IsNull() && !planObj.IsUnknown() {
				planParsed = TFModelObject[T](ctx, diags, planObj)
			}
			if diags.HasError() {
				return nil
			}
			diffs = append(diffs, DiffTPF[T]{
				Path:          p,
				State:         stateParsed,
				Config:        configParsed,
				Plan:          planParsed,
				PlanUnknown:   planObj.IsUnknown(),
				ConfigUnknown: configObj.IsUnknown(),
			})
		}
	}
	return diffs
}
