package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/utility"
)

// DBVersion* returns a plan modifier that applies if there is a known planned value AND user has this defined
// in the configuration.
// This modifier is implemented to maintain backwards compatibility for attributes like mongo_db_major_version
// which allowed user to configure database version in the format like "6.0" or "6" interchangeably and
// formatted the version to "6.0" to persist in the state (Plugin SDKv2's StateFunc). The Terraform plugin framework
// no longer allows the state value of an attribute to be different than the user configured value (if the user does configure a value).
func DBVersion() planmodifier.String {
	return dbVersionModifier{}
}

// useStateForUnknownModifier implements the plan modifier.
type dbVersionModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m dbVersionModifier) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m dbVersionModifier) MarkdownDescription(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// PlanModifyInt64 implements the plan modification logic.
func (m dbVersionModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// if there is a known planned value AND user has this defined in the config, then plan should always show formatted version
	if !req.ConfigValue.IsNull() {
		resp.PlanValue = types.StringValue(utility.FormatMongoDBMajorVersion(req.PlanValue.ValueString()))
		return
	}

	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue = req.StateValue
}
