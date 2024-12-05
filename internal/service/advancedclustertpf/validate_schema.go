package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var defaultMongoDBMajorVersion = "8.0"

func PlanMustUseMongoDBVersion(version float64, mustBeLower bool) FailOnIncompatibleMongoDBVersion {
	return FailOnIncompatibleMongoDBVersion{
		Version:     version,
		MustBeLower: mustBeLower,
	}
}

type FailOnIncompatibleMongoDBVersion struct {
	Version     float64
	MustBeLower bool
}

func (v FailOnIncompatibleMongoDBVersion) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v FailOnIncompatibleMongoDBVersion) MarkdownDescription(_ context.Context) string {
	if v.MustBeLower {
		return fmt.Sprintf("can only be configured if the mongo_db_major_version is lower than %.1f", v.Version)
	}
	return fmt.Sprintf("can only be configured if the mongo_db_major_version is %.1f or higher", v.Version)
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
	isHigherOrEqual := IsMajorVersionHigherGreaterOrEqual(&mongoDbMajorVersionString, v.Version)
	if isHigherOrEqual == nil {
		diags.AddWarning("Unable to parse mongo_db_major_version", "")
		return
	}
	if v.MustBeLower && *isHigherOrEqual || !v.MustBeLower && !*isHigherOrEqual {
		diags.AddError(fmt.Sprintf("`%s` %s", validationPath, v.Description(ctx)), "")
	}
}
