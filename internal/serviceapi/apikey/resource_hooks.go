package apikey

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

// Manual-parity custom code for mongodbatlas_api_key (lift-and-shift from the handwritten SDKv2 resource).

type TFExpandedModel struct {
	// Crafted id of the manual resource, kept so migrated state and references to .id keep working.
	// Injected into the response body by the hooks below under a provider-only key.
	ID types.String `tfsdk:"id" apiname:"tfStateId" autogen:"omitjson"`
}

const stateIDKey = "tfStateId"

var legacyImportIDRegex = regexp.MustCompile(`^([0-9a-f]{24})-([0-9a-f]{24})$`)

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
	return baseSchema
}

func (r *rs) PostCreateAPICall(result autogen.APICallResult) autogen.APICallResult {
	return adaptResponse(result, "")
}

func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	// The API returns 400 API_KEY_NOT_FOUND instead of 404 for a deleted key; the manual resource
	// removed it from state in both cases.
	if validate.StatusBadRequest(result.Resp) {
		result.Err = fmt.Errorf("%w: %v", autogen.ErrNotFound, result.Err)
		return result
	}
	orgID := ""
	if model, ok := req.State.(*TFModel); ok {
		orgID = model.OrgId.ValueString()
	}
	return adaptResponse(result, orgID)
}

func (r *rs) PostUpdateAPICall(result autogen.APICallResult) autogen.APICallResult {
	return adaptResponse(result, "")
}

// PreImport keeps accepting the manual import format {org_id}-{api_key_id}; both parts are
// 24-hex ids, so the dash is unambiguous.
func (r *rs) PreImport(id string) (string, error) {
	if m := legacyImportIDRegex.FindStringSubmatch(id); m != nil {
		return m[1] + "/" + m[2], nil
	}
	return id, nil
}

// adaptResponse flattens the response role objects into the ORG_ role names of the key's
// organization (same as the manual flattenOrgAPIKeyRoles) and injects the crafted id.
// When orgID is empty (create/update responses carry no org id), it is taken from the org roles.
func adaptResponse(result autogen.APICallResult, orgID string) autogen.APICallResult {
	if result.Err != nil || autogen.IsEmptyJSON(result.Body) {
		return result
	}
	var obj map[string]any
	if err := autogen.Decode(result.Body, &obj); err != nil {
		return result
	}
	roles, _ := obj["roles"].([]any)
	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleObj, _ := role.(map[string]any)
		name, _ := roleObj["roleName"].(string)
		roleOrgID, _ := roleObj["orgId"].(string)
		if !strings.HasPrefix(name, "ORG_") {
			continue
		}
		if orgID == "" {
			orgID = roleOrgID
		}
		if roleOrgID == orgID {
			roleNames = append(roleNames, name)
		}
	}
	obj["roles"] = roleNames
	if apiKeyID, _ := obj["id"].(string); apiKeyID != "" && orgID != "" {
		obj[stateIDKey] = conversion.EncodeStateID(map[string]string{
			"org_id":     orgID,
			"api_key_id": apiKeyID,
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
