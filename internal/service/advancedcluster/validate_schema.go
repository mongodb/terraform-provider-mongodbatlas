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

const (
	clusterTypeSharded      = "SHARDED"
	clusterTypeGeosharded   = "GEOSHARDED"
	databaseEditionInfinite = "INFINITE"
)

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
	regionConfigs := newRegionConfig(ctx, req.ConfigValue, diags)
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

// UseEffectiveFieldsValidator validates that use_effective_fields is not set for Flex or Tenant clusters
type UseEffectiveFieldsValidator struct{}

func (v UseEffectiveFieldsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v UseEffectiveFieldsValidator) MarkdownDescription(_ context.Context) string {
	return "use_effective_fields cannot be set for Flex or Tenant clusters."
}

func (v UseEffectiveFieldsValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	diags := &resp.Diagnostics
	var replicationSpecsList types.List
	diags.Append(req.Config.GetAttribute(ctx, path.Root("replication_specs"), &replicationSpecsList)...)
	if diags.HasError() {
		return
	}
	replicationSpecs := newReplicationSpec(ctx, replicationSpecsList, &resp.Diagnostics)
	if diags.HasError() || replicationSpecs == nil {
		return
	}
	if isFlex(replicationSpecs) || isTenant(replicationSpecs) {
		diags.AddAttributeError(req.Path, "Invalid Attribute Configuration",
			"use_effective_fields cannot be set for Flex or Tenant clusters, it is only supported for dedicated clusters.")
	}
}

// InfiniteDatabaseEditionValidator rejects the INFINITE database edition for SHARDED and GEOSHARDED
// clusters at plan time. It only sees config, so it cannot cover values that resolve at apply time or
// updates that omit database_edition; Create and Update gate those paths separately.
type InfiniteDatabaseEditionValidator struct{}

func (v InfiniteDatabaseEditionValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v InfiniteDatabaseEditionValidator) MarkdownDescription(_ context.Context) string {
	return errorInfiniteShardedEdition
}

func (v InfiniteDatabaseEditionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var databaseEdition types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("database_edition"), &databaseEdition)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if isDatabaseEditionInfiniteSharded(req.ConfigValue.ValueString(), databaseEdition.ValueString(), "") {
		resp.Diagnostics.AddAttributeError(req.Path, errorInvalidAttributeConfiguration, errorInfiniteShardedEdition)
	}
}

// isDatabaseEditionInfiniteSharded reports whether the cluster combines a sharded cluster type with the
// INFINITE database edition, which Atlas does not support yet. The configured database_edition takes
// precedence; effective_database_edition from state covers updates that omit database_edition from config.
// Callers pass "" for a value they don't have; "" never matches INFINITE, so those cases are not gated.
// Remove this gate (and the call sites in Create and Update) once Atlas supports these topologies.
func isDatabaseEditionInfiniteSharded(clusterType, databaseEdition, effectiveDatabaseEdition string) bool {
	if clusterType != clusterTypeSharded && clusterType != clusterTypeGeosharded {
		return false
	}
	if databaseEdition != "" {
		return databaseEdition == databaseEditionInfinite
	}
	return effectiveDatabaseEdition == databaseEditionInfinite
}
