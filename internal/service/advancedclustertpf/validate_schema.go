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

func (v FailOnIncompatibleMongoDBVersion) Description(_ context.Context) string {
	if v.MustBeLower {
		return fmt.Sprintf("can only be configured if the mongo_db_major_version is %.1f or lower", v.Version)
	}
	return fmt.Sprintf("can only be configured if the mongo_db_major_version is %.1f or higher", v.Version)
}

func (v FailOnIncompatibleMongoDBVersion) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v FailOnIncompatibleMongoDBVersion) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	performValidation(ctx, &req.Plan, &resp.Diagnostics, v, req.Path)
}

func (v FailOnIncompatibleMongoDBVersion) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	performValidation(ctx, &req.Plan, &resp.Diagnostics, v, req.Path)
}

func performValidation(ctx context.Context, plan *tfsdk.Plan, diags *diag.Diagnostics, v FailOnIncompatibleMongoDBVersion, validationPath path.Path) {
	var mongoDbMajorVersion types.String
	diags.Append(plan.GetAttribute(ctx, path.Root("mongo_db_major_version"), &mongoDbMajorVersion)...)
	if diags.HasError() {
		return
	}
	var mongoDbMajorVersionString string
	if mongoDbMajorVersion.IsUnknown() || mongoDbMajorVersion.IsNull() {
		mongoDbMajorVersionString = defaultMongoDBMajorVersion
	} else {
		mongoDbMajorVersionString = mongoDbMajorVersion.ValueString()
	}
	var validationOk bool
	if v.MustBeLower {
		validationOk = !IsMajorVersionHigherGreaterOrEqual(&mongoDbMajorVersionString, v.Version)
	} else {
		validationOk = IsMajorVersionHigherGreaterOrEqual(&mongoDbMajorVersionString, v.Version)
	}
	if !validationOk {
		diags.AddError(fmt.Sprintf("`%s` %s", validationPath, v.Description(ctx)), "")
	}
}
