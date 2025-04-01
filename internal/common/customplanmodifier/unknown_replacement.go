package customplanmodifier

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

// NewUnknownReplacements creates a new UnknownReplacements instance. ResourceInfo is a struct for storing custom resource specific data. For example, `advanced_cluster` ResourceInfo will differ from `search_deployment` or `project` ResourceInfo
func NewUnknownReplacements[ResourceInfo any](ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, diags *diag.Diagnostics, schema conversion.TPFSchema, info ResourceInfo) *UnknownReplacements[ResourceInfo] {
	differ := NewPlanModifyDiffer(ctx, state, plan, diags, schema)
	return &UnknownReplacements[ResourceInfo]{
		Differ:       differ,
		Info:         info,
		Replacements: make(map[string]UnknownReplacementCall[ResourceInfo]),
	}
}

type UnknownReplacementCall[ResourceInfo any] func(ctx context.Context, stateValue attr.Value, req *UnknownReplacementRequest[ResourceInfo]) attr.Value
type AddKeepUnknownsCall[ResourceInfo any] func(ctx context.Context, stateValue attr.Value, req *UnknownReplacementRequest[ResourceInfo]) []string

type UnknownReplacements[ResourceInfo any] struct {
	Differ       *PlanModifyDiffer
	Replacements map[string]UnknownReplacementCall[ResourceInfo]
	Info         ResourceInfo

	keepUnknownAttributeNames []string // todo: Support validating values when adding attributes
	keepUnknownsExtraCalls    []func(ctx context.Context, stateValue attr.Value, req *UnknownReplacementRequest[ResourceInfo]) []string
}

// AddReplacement call will only be used if the attribute is Unknown in the plan. Only valid for `computed` attributes.
func (u *UnknownReplacements[ResourceInfo]) AddReplacement(name string, call UnknownReplacementCall[ResourceInfo]) {
	// todo: Validate the name exists in the schema and that the attribute is marked with `computed` CLOUDP-309460
	_, found := u.Replacements[name]
	if found {
		panic(fmt.Sprintf("Replacement already exists for %s", name))
	}
	u.Replacements[name] = call
}

// AddKeepUnknownAlways adds the attribute name to the list of attributes that should always keep unknown values. For example connection_string or state_name.
func (u *UnknownReplacements[ResourceInfo]) AddKeepUnknownAlways(keepUnknown ...string) {
	u.keepUnknownAttributeNames = append(u.keepUnknownAttributeNames, keepUnknown...)
}

// AddKeepUnknownOnChanges adds the attribute changed and its depending attributes to the list of attributes that should keep unknown values.
// However, it does not infer dependencies. For example: instance_size --> disk_size_gb, and disk_gb --> disk_iops, doesn't mean instance_size --> disk_iops.
func (u *UnknownReplacements[ResourceInfo]) AddKeepUnknownOnChanges(attributeAffectedMapping map[string][]string) {
	u.keepUnknownAttributeNames = append(u.keepUnknownAttributeNames, u.Differ.AttributeChanges.KeepUnknown(attributeAffectedMapping)...)
}

// AddKeepUnknownsExtraCall adds a function that returns extra keepUnknown attribute names based on the path/stateValue/req (same arguments as the replacer function).
func (u *UnknownReplacements[ResourceInfo]) AddKeepUnknownsExtraCall(call AddKeepUnknownsCall[ResourceInfo]) {
	u.keepUnknownsExtraCalls = append(u.keepUnknownsExtraCalls, call)
}

// ApplyReplacements iterates over the unknown values in the plan and applies the replacement function for each unknown value.
// If there is no explicit replacement function, it will use the default replacer that respects the keepUnknown attributes.
// The calls are done top-down, for example replication_specs.*.id before replication_specs.*.region_configs.*.electable_specs
// Same levels are sorted alphabetically, for example ...region_configs.electable_specs before ...region_configs.read_only_specs
func (u *UnknownReplacements[ResourceInfo]) ApplyReplacements(ctx context.Context, diags *diag.Diagnostics) {
	replacedPaths := []path.Path{}
	ancestorHasProcessed := func(p path.Path) bool {
		for _, replacedPath := range replacedPaths {
			if conversion.HasAncestor(p, replacedPath) {
				return true
			}
		}
		return false
	}
	for _, unknown := range u.Differ.Unknowns(ctx, diags) {
		strPath := unknown.StrPath
		replacer, ok := u.Replacements[unknown.AttributeName]
		if !ok {
			replacer = u.defaultReplacer
		}
		if ancestorHasProcessed(unknown.Path) {
			continue
		}
		replacedPaths = append(replacedPaths, unknown.Path)
		req := &UnknownReplacementRequest[ResourceInfo]{
			Info:          u.Info,
			Path:          unknown.Path,
			Differ:        u.Differ,
			Changes:       u.Differ.AttributeChanges,
			Unknown:       unknown.UnknownValue,
			Diags:         diags,
			AttributeName: unknown.AttributeName,
		}
		replacement := replacer(ctx, unknown.StateValue, req)
		if replacement.IsUnknown() {
			tflog.Debug(ctx, fmt.Sprintf("Keeping unknown value in plan @ %s", strPath))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Replacing unknown value in plan @ %s", strPath))
			UpdatePlanValue(ctx, diags, u.Differ, unknown.Path, replacement)
		}
	}
}

func (u *UnknownReplacements[ResourceInfo]) defaultReplacer(ctx context.Context, stateValue attr.Value, req *UnknownReplacementRequest[ResourceInfo]) attr.Value {
	keepUnknowns := slices.Clone(u.keepUnknownAttributeNames)
	for _, call := range u.keepUnknownsExtraCalls {
		keepUnknowns = append(keepUnknowns, call(ctx, stateValue, req)...)
	}
	if slices.Contains(keepUnknowns, req.AttributeName) {
		return req.Unknown
	}
	return stateValue
}

type UnknownReplacementRequest[ResourceInfo any] struct {
	Info          ResourceInfo      // Resource specific info, useful for storing shardingConfigUpgrade or other relevant info.
	Unknown       attr.Value        // The unknown value in the plan, useful for returning if no replacement is found.
	Differ        *PlanModifyDiffer // Used to read the state and plan values.
	Diags         *diag.Diagnostics
	AttributeName string           // The name of the attribute, for example javascript_enabled for advanced_configuration.javascript_enabled
	Path          path.Path        // The full path to the attribute in the plan, for example advanced_configuration.javascript_enabled
	Changes       AttributeChanges // The changes in the plan, useful for checking if a dependent attribute has changed.
}
