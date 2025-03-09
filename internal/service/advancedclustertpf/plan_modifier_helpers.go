package advancedclustertpf

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

func newDiffHelper(ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, schema SimplifiedSchema) *DiffHelper {
	diags := &resp.Diagnostics
	diffStatePlan, err := req.State.Raw.Diff(resp.Plan.Raw)
	if err != nil {
		diags.AddError("Error diffing state and plan", err.Error())
		return nil
	}
	diffStateConfig, err := req.State.Raw.Diff(req.Config.Raw)
	if err != nil {
		diags.AddError("Error diffing state and config", err.Error())
		return nil
	}
	return &DiffHelper{
		req:             req,
		resp:            resp,
		stateConfigDiff: diffStateConfig,
		statePlanDiff:   diffStatePlan,
		schema:          schema,
	}
}

type DiffHelper struct {
	req             *resource.ModifyPlanRequest
	resp            *resource.ModifyPlanResponse
	stateConfigDiff []tftypes.ValueDiff
	statePlanDiff   []tftypes.ValueDiff
	schema          SimplifiedSchema
}

func (d *DiffHelper) AttributeChanges() schemafunc.AttributeChanges {
	return schemafunc.AttributeChanges{} // TODO
}

func (d *DiffHelper) NiceDiff(ctx context.Context, diags *diag.Diagnostics, schema SimplifiedSchema) string {
	diffPaths := make([]string, len(d.stateConfigDiff))
	for i, diff := range d.statePlanDiff {
		p, localDiags := AttributePath(ctx, diff.Path, schema)
		if localDiags.HasError() {
			diags.Append(localDiags...)
			return ""
		}
		if p.String() == "replication_specs[0].region_configs[1]" {
			continue
		}
		diffPaths[i] = p.String()

	}
	sort.Strings(diffPaths)
	return "Differ\n" + strings.Join(diffPaths, "\n")
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

func keepUnknownCall(aPath *tftypes.AttributePath, keepUnknown []string) bool {
	for _, step := range aPath.Steps() {
		if aName, ok := step.(tftypes.AttributeName); ok {
			if slices.Contains(keepUnknown, string(aName)) {
				return true
			}
		}
	}
	return false
}
func hasPrefix(p path.Path, prefix path.Path) bool {
	prefixString := prefix.String()
	pString := p.String()
	return strings.HasPrefix(pString, prefixString)
}

func UseStateForUnknown(ctx context.Context, diags *diag.Diagnostics, d *DiffHelper, schema SimplifiedSchema, keepUnknown []string, prefix path.Path) {
	for _, diff := range d.statePlanDiff {
		if diff.Value2 != nil && !diff.Value2.IsKnown() {
			stateValue, tpfPath := AttributePathValue(ctx, diags, diff.Path, d.req.State, schema)
			if !hasPrefix(tpfPath, prefix) || stateValue == nil {
				continue
			}
			if keepUnknownCall(diff.Path, keepUnknown) {
				tflog.Info(ctx, fmt.Sprintf("Keeping unknown value in plan @ %s", tpfPath.String()))
				unknownValue := asUnknownValue(ctx, stateValue)
				UpdatePlanValue(ctx, diags, d, tpfPath, unknownValue)
			} else {
				tflog.Info(ctx, fmt.Sprintf("Replacing unknown value in plan @ %s", tpfPath.String()))
				UpdatePlanValue(ctx, diags, d, tpfPath, stateValue)
			}
		}
	}
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
