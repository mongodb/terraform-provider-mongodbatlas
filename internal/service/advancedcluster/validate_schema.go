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

var (
	changeStreamOptionsVersionGate = PlanMustUseMongoDBVersion(7.0, EqualOrHigher)
	defaultMaxTimeMSVersionGate    = PlanMustUseMongoDBVersion(8.0, EqualOrHigher)

	// advancedConfigVersionGates lists the version-gated advanced_configuration attributes. The schema plan
	// modifiers and ValidateMongoDBVersionAtApply both use these gates, so a new gate must be added here too.
	advancedConfigVersionGates = []struct {
		attrName string
		gate     FailOnIncompatibleMongoDBVersion
	}{
		{"change_stream_options_pre_and_post_images_expire_after_seconds", changeStreamOptionsVersionGate},
		{"default_max_time_ms", defaultMaxTimeMSVersionGate},
	}
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
	performValidation(ctx, &req.State, &req.Config, &resp.Diagnostics, v, req.Path)
}

func (v FailOnIncompatibleMongoDBVersion) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	performValidation(ctx, &req.State, &req.Config, &resp.Diagnostics, v, req.Path)
}

// performValidation resolves mongo_db_major_version from config, falling back to state only when the
// attribute is not configured (then it cannot change in this apply). If the version cannot be resolved,
// the check is skipped and Atlas validates at apply.
func performValidation(ctx context.Context, state *tfsdk.State, config *tfsdk.Config, diags *diag.Diagnostics, v FailOnIncompatibleMongoDBVersion, validationPath path.Path) {
	var mongoDbMajorVersionConfig types.String
	var mongoDbMajorVersionState types.String
	diags.Append(config.GetAttribute(ctx, path.Root("mongo_db_major_version"), &mongoDbMajorVersionConfig)...)
	diags.Append(state.GetAttribute(ctx, path.Root("mongo_db_major_version"), &mongoDbMajorVersionState)...)
	if diags.HasError() {
		return
	}
	if mongoDbMajorVersionConfig.IsUnknown() {
		return // set through a reference not known until apply, the version may differ from state
	}
	mongoDbMajorVersionString := mongoDbMajorVersionConfig.ValueString()
	if mongoDbMajorVersionString == "" {
		mongoDbMajorVersionString = mongoDbMajorVersionState.ValueString()
	}
	if mongoDbMajorVersionString == "" {
		return // version unknown (e.g. create without mongo_db_major_version): let Atlas validate at apply
	}
	checkMongoDBVersion(ctx, diags, v, mongoDbMajorVersionString, validationPath)
}

func checkMongoDBVersion(ctx context.Context, diags *diag.Diagnostics, v FailOnIncompatibleMongoDBVersion, version string, validationPath path.Path) {
	isCompatible := MajorVersionCompatible(&version, v.Version, v.Operator)
	if isCompatible == nil {
		diags.AddWarning("Unable to parse mongo_db_major_version", "")
		return
	}
	if !*isCompatible {
		diags.AddError(fmt.Sprintf("%#q %s", validationPath, v.Description(ctx)), "")
	}
}

// ValidateMongoDBVersionAtApply checks the version-gated advanced_configuration attributes in Create and Update,
// before any API call. The plan modifier skips the check when mongo_db_major_version is set through a reference
// unknown at plan time. Without this check, Update would apply the version change and then fail on the advanced
// configuration call, leaving the cluster changed without writing the new state. At apply the reference is
// resolved, so the planned version is used, falling back to state when the attribute is not configured.
// If neither has a version (create with the Atlas default), the check is skipped and Atlas validates.
func ValidateMongoDBVersionAtApply(ctx context.Context, diags *diag.Diagnostics, config *tfsdk.Config, plannedVersion, stateVersion types.String) {
	version := plannedVersion.ValueString() // "" when null or unknown
	if version == "" {
		version = stateVersion.ValueString()
	}
	if version == "" {
		return
	}
	var advancedConfig types.Object
	diags.Append(config.GetAttribute(ctx, path.Root("advanced_configuration"), &advancedConfig)...)
	if diags.HasError() || advancedConfig.IsNull() || advancedConfig.IsUnknown() {
		return
	}
	attrs := advancedConfig.Attributes()
	for _, g := range advancedConfigVersionGates {
		if value, ok := attrs[g.attrName]; !ok || value.IsNull() || value.IsUnknown() {
			continue
		}
		checkMongoDBVersion(ctx, diags, g.gate, version, path.Root("advanced_configuration").AtName(g.attrName))
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
