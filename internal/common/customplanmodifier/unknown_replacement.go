package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewUnknownReplacements[ResourceInfo any](ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, schema conversion.TPFSchema, info ResourceInfo) *UnknownReplacments[ResourceInfo] {
	return &UnknownReplacments[ResourceInfo]{
		Differ: NewPlanModifyDiffer(ctx, req, resp, schema),
		Info:   info,
	}
}

type UnknownReplacementCall[ResourceInfo any] func(ctx context.Context, stateValue ParsedAttrValue, req *UnknownReplacementRequest[ResourceInfo]) attr.Value

type UnknownReplacments[ResourceInfo any] struct {
	Differ       *PlanModifyDiffer
	Replacements map[string]UnknownReplacementCall[ResourceInfo]
	Info         ResourceInfo
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

func (p *ParsedAttrValue) CreateUnknown(ctx context.Context) attr.Value {
	return conversion.AsUnknownValue(ctx, p.Value)
}

type UnknownReplacementRequest[ResourceInfo any] struct {
	Info    ResourceInfo
	Unknown attr.Value
	Differ  *PlanModifyDiffer
	Path    path.Path
	Changes AttributeChanges
}

func (u *UnknownReplacments[ResourceInfo]) AddReplacement(name string, call UnknownReplacementCall[ResourceInfo]) {
	// todo: Validate the name in the schema
	// todo: Validate the name is not already in the replacements
	u.Replacements[name] = call
}

func (u *UnknownReplacments[ResourceInfo]) ApplyReplacments(ctx context.Context, diags *diag.Diagnostics) {
	for strPath, unknown := range u.Differ.Unknowns(ctx, diags) {
		replacer, ok := u.Replacements[unknown.AttributeName]
		if !ok {
			continue
		}
		req := &UnknownReplacementRequest[ResourceInfo]{
			Info:    u.Info,
			Path:    unknown.Path,
			Differ:  u.Differ,
			Changes: u.Differ.AttributeChanges,
			Unknown: unknown.UnknownValue,
		}
		response := replacer(ctx, ParsedAttrValue{Value: unknown.StateValue}, req)
		if response.IsUnknown() {
			tflog.Info(ctx, fmt.Sprintf("Keeping unknown value in plan @ %s", strPath))
		} else {
			tflog.Info(ctx, fmt.Sprintf("Replacing unknown value in plan @ %s", strPath))
			UpdatePlanValue(ctx, diags, u.Differ, unknown.Path, response)
		}
	}
}
