package advancedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var defaultMongoDBMajorVersion = "8.0"

func PlanMustUseMongoDBVersion(version float64, operator MajorVersionOperator) FailOnIncompatibleMongoDBVersion {
	return FailOnIncompatibleMongoDBVersion{
		Version:  version,
		Operator: operator,
	}
}

type FailOnIncompatibleMongoDBVersion struct {
	Version  float64
	Operator MajorVersionOperator
}

func (v FailOnIncompatibleMongoDBVersion) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v FailOnIncompatibleMongoDBVersion) MarkdownDescription(_ context.Context) string {
	switch v.Operator {
	case EqualOrHigher:
		return fmt.Sprintf("can only be configured if the mongo_db_major_version is %.1f or higher", v.Version)
	case Higher:
		return fmt.Sprintf("can only be configured if the mongo_db_major_version is higher than %.1f", v.Version)
	case EqualOrLower:
		return fmt.Sprintf("can only be configured if the mongo_db_major_version is %.1f or lower", v.Version)
	case Lower:
		return fmt.Sprintf("can only be configured if the mongo_db_major_version is lower than %.1f", v.Version)
	default:
		return "unknown operator used"
	}
}

func (v FailOnIncompatibleMongoDBVersion) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	performValidation(ctx, &req.State, &req.Plan, &resp.Diagnostics, v, req.Path)
}

func (v FailOnIncompatibleMongoDBVersion) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	performValidation(ctx, &req.State, &req.Plan, &resp.Diagnostics, v, req.Path)
}

func performValidation(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, diags *diag.Diagnostics, v FailOnIncompatibleMongoDBVersion, validationPath path.Path) {
	var mongoDbMajorVersion types.String
	var mongoDbMajorVersionState types.String
	diags.Append(plan.GetAttribute(ctx, path.Root("mongo_db_major_version"), &mongoDbMajorVersion)...)
	diags.Append(state.GetAttribute(ctx, path.Root("mongo_db_major_version"), &mongoDbMajorVersionState)...)
	if diags.HasError() {
		return
	}
	mongoDbMajorVersionString := mongoDbMajorVersion.ValueString()
	if mongoDbMajorVersionString == "" {
		mongoDbMajorVersionString = mongoDbMajorVersionState.ValueString()
	}
	if mongoDbMajorVersionString == "" {
		mongoDbMajorVersionString = defaultMongoDBMajorVersion
	}
	isCompatible := MajorVersionCompatible(&mongoDbMajorVersionString, v.Version, v.Operator)
	if isCompatible == nil {
		diags.AddWarning("Unable to parse mongo_db_major_version", "")
		return
	}
	if !*isCompatible {
		diags.AddError(fmt.Sprintf("`%s` %s", validationPath, v.Description(ctx)), "")
	}
}

type RegionSpecPriorityOrderDecreasingValidator struct{}

func (v RegionSpecPriorityOrderDecreasingValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}
func (v RegionSpecPriorityOrderDecreasingValidator) MarkdownDescription(_ context.Context) string {
	return "must be a list with priority in descending order"
}
func (v RegionSpecPriorityOrderDecreasingValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	diags := &resp.Diagnostics
	regionConfigs := newCloudRegionConfig20240805(ctx, req.ConfigValue, diags)
	if diags.HasError() || regionConfigs == nil {
		return
	}
	configs := *regionConfigs
	for i := range len(configs) - 1 {
		if configs[i].GetPriority() < configs[i+1].GetPriority() {
			diags.AddError(errorRegionPriorities, fmt.Sprintf("priority value at index %d is %d and priority value at index %d is %d", i, configs[i].GetPriority(), i+1, configs[i+1].GetPriority()))
		}
	}
}
