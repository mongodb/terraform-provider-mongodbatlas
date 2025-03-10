package customplanmodifier

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
)

func NewPlanModifyDiffer(ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, schema conversion.TPFSchema) *PlanModifyDiffer {
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
	schema           conversion.TPFSchema
	AttributeChanges *AttributeChanges
	req              *resource.ModifyPlanRequest
	resp             *resource.ModifyPlanResponse
	stateConfigDiff  []tftypes.ValueDiff
	statePlanDiff    []tftypes.ValueDiff
	PlanFullyKnown   bool
}

func (d *PlanModifyDiffer) ParentRemoved(p path.Path) bool {
	for {
		parent := p.ParentPath()
		if parent.Equal(path.Empty()) {
			return false
		}
		if slices.Contains(*d.AttributeChanges, conversion.AsRemovedIndex(parent)) {
			return true
		}
		p = parent
	}
}

func (d *PlanModifyDiffer) Diff(ctx context.Context, diags *diag.Diagnostics, schema conversion.TPFSchema, isConfig bool) string {
	diffList := d.statePlanDiff
	if isConfig {
		diffList = d.stateConfigDiff
	}
	diffPaths := make([]string, len(diffList))
	for i, diff := range diffList {
		p, localDiags := conversion.AttributePath(ctx, diff.Path, schema)
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
		stateValue, tpfPath := conversion.AttributePathValue(ctx, diags, diff.Path, d.req.State, schema)
		if !conversion.HasPrefix(tpfPath, prefix) || stateValue == nil || conversion.IsAttributeValueOnly(tpfPath) {
			continue
		}
		if d.ParentRemoved(tpfPath) {
			continue
		}
		planValue, _ := conversion.AttributePathValue(ctx, diags, diff.Path, d.req.Plan, schema)
		if planValue == nil || !planValue.IsUnknown() {
			continue
		}
		if keepUnknownCall(diff.Path, keepUnknown) {
			tflog.Info(ctx, fmt.Sprintf("Keeping unknown value in plan @ %s", tpfPath.String()))
			unknownValue := conversion.AsUnknownValue(ctx, stateValue)
			UpdatePlanValue(ctx, diags, d, tpfPath, unknownValue)
		} else {
			tflog.Info(ctx, fmt.Sprintf("Replacing unknown value in plan @ %s", tpfPath.String()))
			UpdatePlanValue(ctx, diags, d, tpfPath, stateValue)
			d.ensureKeepUnknownRespected(ctx, diags, tpfPath, stateValue, keepUnknown)
		}
	}
}

func (d *PlanModifyDiffer) ensureKeepUnknownRespected(ctx context.Context, diags *diag.Diagnostics, tpfPath path.Path, value attr.Value, keepUnknown []string) {
	valueObject, ok := value.(types.Object)
	if value.IsNull() || value.IsUnknown() || !ok {
		return
	}
	for key, childValue := range valueObject.Attributes() {
		if slices.Contains(keepUnknown, key) && !childValue.IsUnknown() {
			childPath := tpfPath.AtName(key)
			tflog.Info(ctx, fmt.Sprintf("Keeping unknown value in plan @ %s", childPath.String()))
			unknownValue := conversion.AsUnknownValue(ctx, childValue)
			UpdatePlanValue(ctx, diags, d, childPath, unknownValue)
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

func readSrcStructValue[T any](ctx context.Context, src conversion.TPFSrc, p path.Path, diags *diag.Diagnostics) *T {
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

func findChanges(ctx context.Context, diff []tftypes.ValueDiff, diags *diag.Diagnostics, schema conversion.TPFSchema) AttributeChanges {
	changes := map[string]bool{}
	addChangeAndParentChanges := func(change string) {
		changes[change] = true
		parts := strings.Split(change, ".")
		for i := range parts[:len(parts)-1] {
			changes[strings.Join(parts[:len(parts)-1-i], ".")] = true
		}
	}
	for _, d := range diff {
		p, localDiags := conversion.AttributePath(ctx, d.Path, schema)
		if conversion.IsListIndex(p) {
			if d.Value1 == nil {
				addChangeAndParentChanges(conversion.AsAddedIndex(p))
			}
			if d.Value2 == nil {
				addChangeAndParentChanges(conversion.AsRemovedIndex(p))
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
		p, localDiags := conversion.AttributePath(ctx, diff.Path, d.schema)
		var parentMatch bool
		if checkNestedAttributes {
			parent := p.ParentPath()
			if conversion.AttributeNameEquals(parent, name) {
				if _, ok := foundParentPaths[parent.String()]; ok {
					continue // parent already used
				}
				foundParentPaths[parent.String()] = true
				p = parent
				parentMatch = true
			}
		}
		// Never show diff if the parent is removed, for example region config
		if !d.ParentRemoved(p) && (parentMatch || conversion.AttributeNameEquals(p, name)) {
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
