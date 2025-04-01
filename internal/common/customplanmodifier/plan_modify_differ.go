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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewPlanModifyDiffer(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, diags *diag.Diagnostics, schema conversion.TPFSchema) *PlanModifyDiffer {
	diffStatePlan, err := state.Raw.Diff(plan.Raw)
	if err != nil {
		diags.AddError("Error diffing state and plan", err.Error())
		return nil
	}

	attributeChanges := findChanges(ctx, diffStatePlan, diags, schema)
	tflog.Debug(ctx, fmt.Sprintf("Attribute changes: %s\n", strings.Join(attributeChanges, "\n")))
	return &PlanModifyDiffer{
		statePlanDiff:    diffStatePlan,
		schema:           schema,
		state:            state,
		plan:             plan,
		AttributeChanges: attributeChanges,
	}
}

type PlanModifyDiffer struct {
	schema           conversion.TPFSchema
	AttributeChanges AttributeChanges
	state            *tfsdk.State
	plan             *tfsdk.Plan
	statePlanDiff    []tftypes.ValueDiff
}

type UnknownInfo struct {
	StateValue    attr.Value
	UnknownValue  attr.Value
	AttributeName string
	StrPath       string
	Path          path.Path
}

func (d *PlanModifyDiffer) Unknowns(ctx context.Context, diags *diag.Diagnostics) []UnknownInfo {
	unknowns := []UnknownInfo{}
	schema := d.schema
	for _, diff := range d.statePlanDiff {
		stateValue, tpfPath := conversion.AttributePathValue(ctx, diags, diff.Path, d.state, schema)
		strPath := tpfPath.String()
		planValue, _ := conversion.AttributePathValue(ctx, diags, diff.Path, d.plan, schema)
		if planValue == nil || !planValue.IsUnknown() {
			continue
		}
		unknowns = append(unknowns, UnknownInfo{
			Path:          tpfPath,
			StrPath:       strPath,
			StateValue:    stateValue,
			UnknownValue:  planValue,
			AttributeName: conversion.AttributeName(tpfPath),
		})
	}
	slices.SortFunc(unknowns, func(i, j UnknownInfo) int {
		return strings.Compare(i.StrPath, j.StrPath)
	})
	return unknowns
}

func ReadPlanStructValue[T any](ctx context.Context, d *PlanModifyDiffer, p path.Path) *T {
	return readSrcStructValue[T](ctx, d.plan, p)
}

func ReadStateStructValue[T any](ctx context.Context, d *PlanModifyDiffer, p path.Path) *T {
	return readSrcStructValue[T](ctx, d.state, p)
}

func readSrcStructValue[T any](ctx context.Context, src conversion.TPFSrc, p path.Path) *T {
	var obj types.Object
	if localDiags := src.GetAttribute(ctx, p, &obj); len(localDiags) > 0 {
		tflog.Error(ctx, conversion.FormatDiags(&localDiags))
		return nil
	}
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	return conversion.TFModelObject[T](ctx, obj)
}

func ReadPlanStructValues[T any](ctx context.Context, d *PlanModifyDiffer, p path.Path) []T {
	return readSrcStructValues[T](ctx, d.plan, p)
}

func readSrcStructValues[T any](ctx context.Context, src conversion.TPFSrc, p path.Path) []T {
	var objList types.List
	var localDiags diag.Diagnostics
	if localDiags = src.GetAttribute(ctx, p, &objList); len(localDiags) > 0 {
		tflog.Error(ctx, conversion.FormatDiags(&localDiags))
		return nil
	}
	result := conversion.TFModelList[T](ctx, &localDiags, objList)
	if len(localDiags) > 0 {
		tflog.Error(ctx, conversion.FormatDiags(&localDiags))
	}
	return result
}

func UpdatePlanValue(ctx context.Context, diags *diag.Diagnostics, d *PlanModifyDiffer, p path.Path, value attr.Value) {
	diags.Append(d.plan.SetAttribute(ctx, p, value)...)
}

func findChanges(ctx context.Context, diff []tftypes.ValueDiff, diags *diag.Diagnostics, schema conversion.TPFSchema) AttributeChanges {
	changes := make(map[string]bool)
	addChangeAndAncestorChanges := func(change path.Path) {
		changes[change.String()] = true
		for _, p := range conversion.AncestorPaths(change) {
			changes[p.String()] = true
		}
	}
	// avoids adding change for removed region_configs inside a removed replication_specs
	isAncestorRemoved := func(p path.Path) bool {
		for _, a := range conversion.AncestorPaths(p) {
			if conversion.IsListIndex(a) {
				if changes[conversion.AsRemovedIndex(a)] {
					return true
				}
			}
		}
		return false
	}
	for _, d := range diff {
		p, localDiags := conversion.AttributePath(ctx, d.Path, schema)
		if localDiags.HasError() {
			diags.Append(localDiags...)
			continue
		}
		// Two types of changes from a diff
		// 1. It is defined in the plan AND it is Known and not null
		// 2. It is a removed list index. (set or map index we ignore, e.g., replication_specs[0].container_id[-\"AWS:US_EAST_1\"] or advanced_configuration.custom_openssl_cipher_config_tls12[-Value(\"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384\")])
		// If we use schema.SetNestedAttribute we might need to see if those changes are detected
		if d.Value2 != nil && d.Value2.IsKnown() && !d.Value2.IsNull() {
			addChangeAndAncestorChanges(p)
		}
		if conversion.IsListIndex(p) {
			isAdd := d.Value1 == nil
			if isAdd {
				changes[conversion.AsAddedIndex(p)] = true
			}
			isRemove := d.Value2 == nil
			if isRemove && !isAncestorRemoved(p) {
				changes[conversion.AsRemovedIndex(p)] = true
			}
		}
	}
	return slices.Sorted(maps.Keys(changes)) // prettier attribute changes when they are sorted alphabetically
}
