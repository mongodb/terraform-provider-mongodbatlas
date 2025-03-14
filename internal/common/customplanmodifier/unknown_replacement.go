package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

type TPFType interface {
	types.Bool | types.Int64 | types.Int32 | types.Float64 | types.Float32 | types.String | types.List | types.Map | types.Object | types.Set | types.Tuple | types.Dynamic
}

func NewUnknownReplacments[ResourceInfo any](ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, schema conversion.TPFSchema, info ResourceInfo) *UnknownReplacments[ResourceInfo] {
	return &UnknownReplacments[ResourceInfo]{
		Differ: NewPlanModifyDiffer(ctx, req, resp, schema),
	}
}

type UnknownAttributeReplacment[ResourceInfo any, T attr.Value] struct {
	Name string
	Call func(T, AttributeChanges, ResourceInfo, *PlanModifyDiffer) T
}

type internalUnknownAttributeReplacment[ResourceInfo any] struct {
	Name string
	Call func(attr.Value, AttributeChanges, ResourceInfo, *PlanModifyDiffer) attr.Value
}
type UnknownReplacments[ResourceInfo any] struct {
	Differ       *PlanModifyDiffer
	Replacements map[string]internalUnknownAttributeReplacment[ResourceInfo]
	ResourceInfo ResourceInfo
}

func (u *UnknownReplacments[ResourceInfo]) ApplyReplacments(ctx context.Context, diags *diag.Diagnostics) {
	for strPath, unknown := range u.Differ.Unknowns(ctx, diags) {
		replacer, ok := u.Replacements[unknown.AttributeName]
		if !ok {
			continue
		}
		response := replacer.Call(unknown.StateValue, u.Differ.AttributeChanges, u.ResourceInfo, u.Differ)
		if response == nil {
			tflog.Info(ctx, fmt.Sprintf("Keeping unknown value in plan @ %s", strPath))
		} else {
			tflog.Info(ctx, fmt.Sprintf("Replacing unknown value in plan @ %s", strPath))
			UpdatePlanValue(ctx, diags, u.Differ, unknown.Path, response)
		}
	}
}

func AddReplacment[ResourceInfo any, T attr.Value](r *UnknownReplacments[ResourceInfo], name string, call func(T, AttributeChanges, ResourceInfo, *PlanModifyDiffer) T) {
	// todo: Validate the name in the schema
	// todo: Validate the name is not already in the replacements
	r.Replacements[name] = internalUnknownAttributeReplacment[ResourceInfo]{
		Name: name,
		Call: func(a attr.Value, c AttributeChanges, r ResourceInfo, d *PlanModifyDiffer) attr.Value {
			asParsed := a.(T)
			return call(asParsed, c, r, d)
		},
	}
}
