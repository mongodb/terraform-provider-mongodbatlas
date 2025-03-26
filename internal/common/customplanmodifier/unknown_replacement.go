package customplanmodifier

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewUnknownReplacements[ResourceInfo any](ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, schema conversion.TPFSchema, info ResourceInfo) *UnknownReplacements[ResourceInfo] {
	differ := NewPlanModifyDiffer(ctx, req, resp, schema)
	tflog.Debug(ctx, differ.Diff(ctx, &resp.Diagnostics, schema, false))
	return &UnknownReplacements[ResourceInfo]{
		Differ:       differ,
		Info:         info,
		Replacements: make(map[string]UnknownReplacementCall[ResourceInfo]),
	}
}

type UnknownReplacementCall[ResourceInfo any] func(ctx context.Context, stateValue ParsedAttrValue, req *UnknownReplacementRequest[ResourceInfo]) attr.Value

type UnknownReplacements[ResourceInfo any] struct {
	Differ       *PlanModifyDiffer
	Replacements map[string]UnknownReplacementCall[ResourceInfo]
	Info         ResourceInfo

	keepUnknownAttributeNames []string // todo: Support validating values when adding attributes
	keepUnknownsExtraCalls    []func(ctx context.Context, stateValue ParsedAttrValue, req *UnknownReplacementRequest[ResourceInfo]) []string
}

func (u *UnknownReplacements[ResourceInfo]) AddReplacement(name string, call UnknownReplacementCall[ResourceInfo]) {
	// todo: Validate the name exists in the schema
	_, existing := u.Replacements[name]
	if existing {
		panic(fmt.Sprintf("Replacement already exists for %s", name))
	}
	u.Replacements[name] = call
}

func (u *UnknownReplacements[ResourceInfo]) AddKeepUnknownAlways(keepUnknown ...string) {
	u.keepUnknownAttributeNames = append(u.keepUnknownAttributeNames, keepUnknown...)
}

func (u *UnknownReplacements[ResourceInfo]) AddKeepUnknownOnChanges(attributeEffectedMapping map[string][]string) {
	u.keepUnknownAttributeNames = append(u.keepUnknownAttributeNames, u.Differ.AttributeChanges.KeepUnknown(attributeEffectedMapping)...)
}

func (u *UnknownReplacements[ResourceInfo]) AddKeepUnknownsExtraCall(call func(ctx context.Context, stateValue ParsedAttrValue, req *UnknownReplacementRequest[ResourceInfo]) []string) {
	u.keepUnknownsExtraCalls = append(u.keepUnknownsExtraCalls, call)
}

func (u *UnknownReplacements[ResourceInfo]) ApplyReplacements(ctx context.Context, diags *diag.Diagnostics) {
	for strPath, unknown := range u.Differ.Unknowns(ctx, diags) {
		replacer, ok := u.Replacements[unknown.AttributeName]
		if !ok {
			replacer = u.defaultReplacer
		}
		req := &UnknownReplacementRequest[ResourceInfo]{
			Info:          u.Info,
			Path:          unknown.Path,
			Differ:        u.Differ,
			Changes:       u.Differ.AttributeChanges,
			Unknown:       unknown.UnknownValue,
			Diags:         diags,
			AttributeName: unknown.AttributeName,
		}
		replacement := replacer(ctx, ParsedAttrValue{Value: unknown.StateValue}, req)
		if replacement.IsUnknown() {
			tflog.Debug(ctx, fmt.Sprintf("Keeping unknown value in plan @ %s", strPath))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Replacing unknown value in plan @ %s", strPath))
			UpdatePlanValue(ctx, diags, u.Differ, unknown.Path, replacement)
		}
	}
}

func (u *UnknownReplacements[ResourceInfo]) defaultReplacer(ctx context.Context, stateValue ParsedAttrValue, req *UnknownReplacementRequest[ResourceInfo]) attr.Value {
	keepUnknowns := slices.Clone(u.keepUnknownAttributeNames)
	for _, call := range u.keepUnknownsExtraCalls {
		keepUnknowns = append(keepUnknowns, call(ctx, stateValue, req)...)
	}
	if slices.Contains(keepUnknowns, req.AttributeName) {
		return req.Unknown
	}
	return stateValue.Value
}

// ParsedAttrValue is a wrapper around attr.Value that provides type-safe accessors to support using the same signature of functions.
type ParsedAttrValue struct {
	Value attr.Value
}

func (p *ParsedAttrValue) AsString() types.String {
	return p.Value.(types.String)
}

func (p *ParsedAttrValue) AsObject() types.Object {
	return p.Value.(types.Object)
}

type UnknownReplacementRequest[ResourceInfo any] struct {
	Info          ResourceInfo
	Unknown       attr.Value
	Differ        *PlanModifyDiffer
	Diags         *diag.Diagnostics
	AttributeName string
	Path          path.Path
	Changes       AttributeChanges
}
