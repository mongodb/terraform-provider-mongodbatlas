package customdbrole

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

// Manual-parity custom code for mongodbatlas_custom_db_role (lift-and-shift from the handwritten SDKv2 resource).

type TFExpandedModel struct {
	// Crafted id of the manual resource, kept so migrated state and references to .id keep working.
	// Injected into the response body by the hooks below under a provider-only key.
	ID types.String `tfsdk:"id" apiname:"tfStateId" autogen:"omitjson"`
}

const stateIDKey = "tfStateId"

var (
	// The project id is always 24 hex characters, so the manual {project_id}-{role_name} format is unambiguous
	// even when the role name contains dashes.
	legacyImportIDRegex = regexp.MustCompile(`^([0-9a-f]{24})-(.+)$`)
	// Create and update responses do not include the project id, it is taken from the request path.
	projectIDPathRegex = regexp.MustCompile(`/groups/([0-9a-f]{24})/customDBRoles/`)
	// Ported as-is from the manual resource: unanchored, so it accepts any value with one word character.
	roleNameRegex = regexp.MustCompile(`[\w-]+`)
)

var _ autogen.ResourceSchemaHook = (*rs)(nil)
var _ autogen.PostCreateAPICallHook = (*rs)(nil)
var _ autogen.PostReadAPICallHook = (*rs)(nil)
var _ autogen.PostUpdateAPICallHook = (*rs)(nil)
var _ autogen.PreImportHook = (*rs)(nil)

func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	baseSchema.Attributes["id"] = schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	// ForceNew in the manual resource: RequiresReplace instead of the generated CreateOnly.
	projectID := baseSchema.Attributes["project_id"].(schema.StringAttribute)
	projectID.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
	baseSchema.Attributes["project_id"] = projectID

	roleName := baseSchema.Attributes["role_name"].(schema.StringAttribute)
	roleName.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
	roleName.Validators = []validator.String{
		stringvalidator.RegexMatches(roleNameRegex, "`role_name` can contain only letters, digits, underscores, and dashes"),
		roleNamePrefixValidator{},
	}
	baseSchema.Attributes["role_name"] = roleName

	// SDKv2 stored unset optional values as "" and false: defaults keep plan, state and migrated state aligned
	// with the always-populated API response.
	actions := baseSchema.Attributes["actions"].(schema.SetNestedAttribute)
	resources := actions.NestedObject.Attributes["resources"].(schema.SetNestedAttribute)
	cluster := resources.NestedObject.Attributes["cluster"].(schema.BoolAttribute)
	cluster.Default = booldefault.StaticBool(false)
	resources.NestedObject.Attributes["cluster"] = cluster
	for _, name := range []string{"collection_name", "database_name"} {
		attr := resources.NestedObject.Attributes[name].(schema.StringAttribute)
		attr.Default = stringdefault.StaticString("")
		resources.NestedObject.Attributes[name] = attr
	}
	actions.NestedObject.Attributes["resources"] = resources
	actions.Computed = true
	actions.PlanModifiers = []planmodifier.Set{emptySetFromSDKv2{}}
	baseSchema.Attributes["actions"] = actions

	inheritedRoles := baseSchema.Attributes["inherited_roles"].(schema.SetNestedAttribute)
	inheritedRoles.Computed = true
	inheritedRoles.PlanModifiers = []planmodifier.Set{emptySetFromSDKv2{}}
	baseSchema.Attributes["inherited_roles"] = inheritedRoles
	return baseSchema
}

// emptySetFromSDKv2 avoids a one-time `[] -> null` diff after migration: SDKv2 stored absent blocks as empty
// sets, while an absent attribute is null. When the config is null it keeps an empty state value and plans null
// otherwise, so the attribute behaves as Optional (Computed is only needed to plan a value different from config).
type emptySetFromSDKv2 struct{}

func (m emptySetFromSDKv2) Description(_ context.Context) string {
	return "Keeps an empty set in state when the attribute is not configured."
}

func (m emptySetFromSDKv2) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m emptySetFromSDKv2) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}
	if !req.StateValue.IsNull() && len(req.StateValue.Elements()) == 0 {
		resp.PlanValue = req.StateValue
		return
	}
	resp.PlanValue = types.SetNull(req.PlanValue.ElementType(ctx))
}

func (r *rs) PostCreateAPICall(result autogen.APICallResult) autogen.APICallResult {
	return adaptResponse(result, projectIDFromRequest(result))
}

func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	projectID := ""
	if model, ok := req.State.(*TFModel); ok {
		projectID = model.ProjectId.ValueString()
	}
	return adaptResponse(result, projectID)
}

func (r *rs) PostUpdateAPICall(result autogen.APICallResult) autogen.APICallResult {
	return adaptResponse(result, projectIDFromRequest(result))
}

// PreImport keeps accepting the manual import format {project_id}-{role_name}.
func (r *rs) PreImport(id string) (string, error) {
	if m := legacyImportIDRegex.FindStringSubmatch(id); m != nil {
		return m[1] + "/" + m[2], nil
	}
	return id, nil
}

func projectIDFromRequest(result autogen.APICallResult) string {
	if result.Resp == nil || result.Resp.Request == nil || result.Resp.Request.URL == nil {
		return ""
	}
	if m := projectIDPathRegex.FindStringSubmatch(result.Resp.Request.URL.Path); m != nil {
		return m[1]
	}
	return ""
}

// adaptResponse fills the action resource fields the response may omit with the values the manual resource
// stored ("" and false) and injects the crafted id.
func adaptResponse(result autogen.APICallResult, projectID string) autogen.APICallResult {
	if result.Err != nil || autogen.IsEmptyJSON(result.Body) {
		return result
	}
	var obj map[string]any
	if err := autogen.Decode(result.Body, &obj); err != nil {
		result.Err = err
		return result
	}
	normalizeActions(obj)
	if roleName, _ := obj["roleName"].(string); roleName != "" && projectID != "" {
		obj[stateIDKey] = conversion.EncodeStateID(map[string]string{
			"project_id": projectID,
			"role_name":  roleName,
		})
	}
	body, err := json.Marshal(obj)
	if err != nil {
		result.Err = err
		return result
	}
	result.Body = body
	return result
}

func normalizeActions(obj map[string]any) {
	actions, _ := obj["actions"].([]any)
	for _, action := range actions {
		actionObj, _ := action.(map[string]any)
		resources, _ := actionObj["resources"].([]any)
		for _, res := range resources {
			resObj, ok := res.(map[string]any)
			if !ok {
				continue
			}
			setIfMissing(resObj, "cluster", false)
			setIfMissing(resObj, "collection", "")
			setIfMissing(resObj, "db", "")
		}
	}
}

func setIfMissing(obj map[string]any, key string, value any) {
	if v, ok := obj[key]; !ok || v == nil {
		obj[key] = value
	}
}

// roleNamePrefixValidator ports the manual prefix check as-is: it rejects names starting with "x-gen"
// while its message says "xgen-".
type roleNamePrefixValidator struct{}

func (v roleNamePrefixValidator) Description(_ context.Context) string {
	return "`role_name` cannot start with 'xgen-'"
}

func (v roleNamePrefixValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v roleNamePrefixValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if strings.HasPrefix(req.ConfigValue.ValueString(), "x-gen") {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid role_name", fmt.Sprintf("%s, got: %s", v.Description(ctx), req.ConfigValue.ValueString()))
	}
}
