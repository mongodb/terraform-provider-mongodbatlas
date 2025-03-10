package advancedclustertpf

import (
	"context"
	"fmt"
	"maps"
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
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

func NewPlanModifyDiffer(ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, schema TPFSchema) *PlanModifyDiffer {
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

	attributeChanges := findChanges(ctx, diffStatePlan, diags, schema)
	tflog.Info(ctx, fmt.Sprintf("Attribute changes: %s\n", strings.Join(attributeChanges, "\n")))
	return &PlanModifyDiffer{
		req:              req,
		resp:             resp,
		stateConfigDiff:  diffStateConfig,
		statePlanDiff:    diffStatePlan,
		schema:           schema,
		AttributeChanges: &attributeChanges,
		PlanFullyKnown:   req.Plan.Raw.IsFullyKnown(),
	}
}

type PlanModifyDiffer struct {
	schema           TPFSchema
	AttributeChanges *schemafunc.AttributeChanges
	req              *resource.ModifyPlanRequest
	resp             *resource.ModifyPlanResponse
	stateConfigDiff  []tftypes.ValueDiff
	statePlanDiff    []tftypes.ValueDiff
	PlanFullyKnown   bool
}

func (d *PlanModifyDiffer) ParentRemoved(p path.Path) bool {
	if !IsListIndex(p.ParentPath()) {
		return false
	}
	parentRemoved := AsRemovedIndex(p.ParentPath())
	return slices.Contains(*d.AttributeChanges, parentRemoved)
}

func (d *PlanModifyDiffer) Diff(ctx context.Context, diags *diag.Diagnostics, schema TPFSchema, isConfig bool) string {
	diffList := d.statePlanDiff
	if isConfig {
		diffList = d.stateConfigDiff
	}
	diffPaths := make([]string, len(diffList))
	for i, diff := range diffList {
		p, localDiags := AttributePath(ctx, diff.Path, schema)
		if localDiags.HasError() {
			diags.Append(localDiags...)
			return ""
		}
		diffPaths[i] = p.String()
	}
	sort.Strings(diffPaths)
	name := "plan"
	if isConfig {
		name = "config"
	}
	return fmt.Sprintf("DifferStateTo%s\n", name) + strings.Join(diffPaths, "\n")
}

func (d *PlanModifyDiffer) UseStateForUnknown(ctx context.Context, diags *diag.Diagnostics, keepUnknown []string, prefix path.Path) {
	// The diff is sorted by the path length, for example read_only_spec is processed before read_only_spec.disk_size_gb
	schema := d.schema
	for _, diff := range d.statePlanDiff {
		stateValue, tpfPath := AttributePathValue(ctx, diags, diff.Path, d.req.State, schema)
		if !HasPrefix(tpfPath, prefix) || stateValue == nil || IsAttributeValueOnly(tpfPath) {
			continue
		}
		planValue, _ := AttributePathValue(ctx, diags, diff.Path, d.req.Plan, schema)
		if planValue == nil {
			continue
		}
		// For nested attributes with unknown values, all their children attributes will be `null` instead of unknown.
		// Therefore, to ensure keepUnknown, force unknown when the responsePlanValue is not unknown.
		if planValue.IsNull() && keepUnknownCall(diff.Path, keepUnknown) {
			responsePlanValue, _ := AttributePathValue(ctx, diags, diff.Path, d.resp.Plan, schema)
			if responsePlanValue != nil && !responsePlanValue.IsUnknown() {
				tflog.Info(ctx, fmt.Sprintf("Force unknown value in plan @ %s", tpfPath.String()))
				unknownValue := asUnknownValue(ctx, stateValue)
				UpdatePlanValue(ctx, diags, d, tpfPath, unknownValue)
			}
			continue
		}
		if !planValue.IsUnknown() {
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

func ReadConfigStructValue[T any](ctx context.Context, diags *diag.Diagnostics, d *PlanModifyDiffer, p path.Path) *T {
	return readSrcStructValue[T](ctx, d.req.Config, p, diags)
}

func ReadPlanStructValue[T any](ctx context.Context, diags *diag.Diagnostics, d *PlanModifyDiffer, p path.Path) *T {
	return readSrcStructValue[T](ctx, d.req.Plan, p, diags)
}

func ReadStateStructValue[T any](ctx context.Context, diags *diag.Diagnostics, d *PlanModifyDiffer, p path.Path) *T {
	return readSrcStructValue[T](ctx, d.req.State, p, diags)
}

func readSrcStructValue[T any](ctx context.Context, src TPFSrc, p path.Path, diags *diag.Diagnostics) *T {
	var obj types.Object
	if localDiags := src.GetAttribute(ctx, p, &obj); localDiags.HasError() {
		diags.Append(localDiags...)
		return nil
	}
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	return conversion.TFModelObject[T](ctx, diags, obj)
}

func UpdatePlanValue[T attr.Value](ctx context.Context, diags *diag.Diagnostics, d *PlanModifyDiffer, p path.Path, value T) {
	if localDiags := d.resp.Plan.SetAttribute(ctx, p, value); localDiags.HasError() {
		diags.Append(localDiags...)
	}
}

type DiffTPF[T any] struct {
	Plan          *T
	State         *T
	Config        *T
	Path          path.Path
	PlanUnknown   bool
	ConfigUnknown bool
}

func (d *DiffTPF[T]) Removed() bool {
	return d.State != nil && d.Config == nil
}

func (d *DiffTPF[T]) Changed() bool {
	return d.State != nil && d.Config != nil
}

func (d *DiffTPF[T]) PlanOrStateValue() *T {
	if d.Plan != nil {
		return d.Plan
	}
	return d.State
}

func findChanges(ctx context.Context, diff []tftypes.ValueDiff, diags *diag.Diagnostics, schema TPFSchema) schemafunc.AttributeChanges {
	changes := map[string]bool{}
	addChangeAndParentChanges := func(change string) {
		changes[change] = true
		parts := strings.Split(change, ".")
		for i := range parts[:len(parts)-1] {
			changes[strings.Join(parts[:len(parts)-1-i], ".")] = true
		}
	}
	for _, d := range diff {
		p, localDiags := AttributePath(ctx, d.Path, schema)
		if IsListIndex(p) {
			if d.Value1 == nil {
				addChangeAndParentChanges(AsAddedIndex(p))
			}
			if d.Value2 == nil {
				addChangeAndParentChanges(AsRemovedIndex(p))
			}
		}
		if d.Value2 != nil && d.Value2.IsKnown() && !d.Value2.IsNull() {
			if localDiags.HasError() {
				diags.Append(localDiags...)
				continue
			}
			addChangeAndParentChanges(p.String())
		}
	}
	return slices.Sorted(maps.Keys(changes))
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

func StateConfigDiffs[T any](ctx context.Context, diags *diag.Diagnostics, d *PlanModifyDiffer, name string, checkNestedAttributes bool) []DiffTPF[T] {
	earlyReturn := func(localDiags diag.Diagnostics) []DiffTPF[T] {
		diags.Append(localDiags...)
		return nil
	}
	var diffs []DiffTPF[T]
	foundParentPaths := map[string]bool{}

	for _, diff := range d.stateConfigDiff {
		p, localDiags := AttributePath(ctx, diff.Path, d.schema)
		var pathMatch bool
		if checkNestedAttributes {
			parent := p.ParentPath()
			if AttributeNameEquals(parent, name) {
				if _, ok := foundParentPaths[parent.String()]; ok {
					continue // parent already used
				}
				foundParentPaths[parent.String()] = true
				p = parent
				pathMatch = true
			}
		}
		// Never show diff if the parent is removed, for exampl region config
		if !d.ParentRemoved(p) && (pathMatch || AttributeNameEquals(p, name)) {
			if localDiags.HasError() {
				return earlyReturn(localDiags)
			}
			var configObj, planObj types.Object
			stateParsed := ReadStateStructValue[T](ctx, diags, d, p)
			if d2 := d.req.Config.GetAttribute(ctx, p, &configObj); d2.HasError() {
				return earlyReturn(d2)
			}
			if d3 := d.req.Plan.GetAttribute(ctx, p, &planObj); d3.HasError() {
				return earlyReturn(d3)
			}
			var configParsed, planParsed *T
			if !configObj.IsNull() && !configObj.IsUnknown() {
				configParsed = conversion.TFModelObject[T](ctx, diags, configObj)
			}
			if !planObj.IsNull() && !planObj.IsUnknown() {
				planParsed = conversion.TFModelObject[T](ctx, diags, planObj)
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
