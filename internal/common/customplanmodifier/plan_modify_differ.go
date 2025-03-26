package customplanmodifier

import (
	"context"
	"fmt"
	"maps"
	"slices"
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
	tflog.Debug(ctx, fmt.Sprintf("Attribute changes: %s\n", strings.Join(attributeChanges, "\n")))
	return &PlanModifyDiffer{
		req:              req,
		resp:             resp,
		stateConfigDiff:  diffStateConfig,
		statePlanDiff:    diffStatePlan,
		schema:           schema,
		AttributeChanges: attributeChanges,
		PlanFullyKnown:   req.Plan.Raw.IsFullyKnown(),
	}
}

type PlanModifyDiffer struct {
	schema           conversion.TPFSchema
	AttributeChanges AttributeChanges
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
		if slices.Contains(d.AttributeChanges, conversion.AsRemovedIndex(parent)) {
			return true
		}
		p = parent
	}
}

type UnknownInfo struct {
	StateValue    attr.Value
	UnknownValue  attr.Value
	AttributeName string
	Path          path.Path
}

func (d *PlanModifyDiffer) Unknowns(ctx context.Context, diags *diag.Diagnostics) map[string]UnknownInfo {
	unknowns := map[string]UnknownInfo{}
	schema := d.schema
	for _, diff := range d.statePlanDiff {
		stateValue, tpfPath := conversion.AttributePathValue(ctx, diags, diff.Path, d.req.State, schema)
		if d.ParentRemoved(tpfPath) {
			continue
		}
		planValue, _ := conversion.AttributePathValue(ctx, diags, diff.Path, d.req.Plan, schema)
		if planValue == nil || !planValue.IsUnknown() {
			continue
		}
		unknowns[tpfPath.String()] = UnknownInfo{
			Path:          tpfPath,
			StateValue:    stateValue,
			UnknownValue:  planValue,
			AttributeName: conversion.AttributeName(tpfPath),
		}
	}
	return unknowns
}

func ReadConfigStructValue[T any](ctx context.Context, d *PlanModifyDiffer, p path.Path) *T {
	return readSrcStructValue[T](ctx, d.req.Config, p)
}

func ReadPlanStructValue[T any](ctx context.Context, d *PlanModifyDiffer, p path.Path) *T {
	return readSrcStructValue[T](ctx, d.req.Plan, p)
}

func ReadStateStructValue[T any](ctx context.Context, d *PlanModifyDiffer, p path.Path) *T {
	return readSrcStructValue[T](ctx, d.req.State, p)
}

func readSrcStructValue[T any](ctx context.Context, src conversion.TPFSrc, p path.Path) *T {
	var obj types.Object
	if localDiags := src.GetAttribute(ctx, p, &obj); localDiags.HasError() {
		return nil
	}
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	return conversion.TFModelObject[T](ctx, obj)
}
func ReadPlanStructValues[T any](ctx context.Context, d *PlanModifyDiffer, p path.Path, diags *diag.Diagnostics) []T {
	return readSrcStructValues[T](ctx, d.req.Plan, p, diags)
}

func readSrcStructValues[T any](ctx context.Context, src conversion.TPFSrc, p path.Path, diags *diag.Diagnostics) []T {
	var objList types.List
	if localDiags := src.GetAttribute(ctx, p, &objList); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return conversion.TFModelList[T](ctx, diags, objList)
}

func UpdatePlanValue(ctx context.Context, diags *diag.Diagnostics, d *PlanModifyDiffer, p path.Path, value attr.Value) {
	diags.Append(d.resp.Plan.SetAttribute(ctx, p, value)...)
}

func findChanges(ctx context.Context, diff []tftypes.ValueDiff, diags *diag.Diagnostics, schema conversion.TPFSchema) AttributeChanges {
	changes := make(map[string]bool)
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
	return slices.Sorted(maps.Keys(changes)) // Ensure changes are sorted to support top-down processing, for example read_only_spec is processed before read_only_spec.disk_size_gb
}
