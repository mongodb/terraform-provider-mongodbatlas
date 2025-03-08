package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type DiffHelper struct {
	req             *resource.ModifyPlanRequest
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

type DiffTPF[T any] struct {
	Path     path.Path
	OldValue *T
	NewValue *T
}

func StateConfigDiffs[T any](ctx context.Context, diags *diag.Diagnostics, d *DiffHelper, name tftypes.AttributeName, schema SimplifiedSchema) []DiffTPF[T] {
	earlyReturn := func(localDiags diag.Diagnostics) []DiffTPF[T] {
		diags.Append(localDiags...)
		return nil
	}
	var diffs []DiffTPF[T]

	for _, diff := range d.stateConfigDiff {
		if diff.Path.LastStep().Equal(name) {
			var stateObj types.Object
			var configObj types.Object
			p, localDiags := AttributePath(ctx, diff.Path, schema)
			if localDiags.HasError() {
				return earlyReturn(localDiags)
			}
			if d1 := d.req.State.GetAttribute(ctx, p, &stateObj); d1.HasError() {
				return earlyReturn(d1)
			}
			if d2 := d.req.Config.GetAttribute(ctx, p, &configObj); d2.HasError() {
				return earlyReturn(d2)
			}
			var configParsed, stateParsed *T
			if !stateObj.IsNull() {
				stateParsed = TFModelObject[T](ctx, diags, stateObj)
			}
			if !configObj.IsNull() {
				configParsed = TFModelObject[T](ctx, diags, configObj)
			}
			if diags.HasError() {
				return nil
			}
			diffs = append(diffs, DiffTPF[T]{Path: p, OldValue: stateParsed, NewValue: configParsed})
		}
	}
	return diffs
}
