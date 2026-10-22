package databaseuser

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

// Manual-parity custom code for mongodbatlas_database_user (spike pilot).
//
// The `validators`/`plan_modifiers` config overrides are not consumed by the generator
// yet, and the write-only password fields / crafted `id` are not in the OpenAPI spec,
// so all of the below lives in hooks (expanded_model: true provides the model slots).

type TFExpandedModel struct {
	// Crafted id reproducing the manual resource format so existing state keeps working
	// (lift-and-shift at the same resource address).
	ID types.String `tfsdk:"id" apiname:"id" autogen:"omitjson"`

	// Write-only password support (manual parity). Marshaled into the request body as
	// `password` via apiname when set; WriteOnly keeps it out of state. Sending an
	// unchanged value on every update is behaviorally equivalent to the manual
	// version-gated flow, except that an out-of-band password change is overwritten
	// on unrelated updates (documented nuance in the spike doc).
	PasswordWo types.String `tfsdk:"password_wo" apiname:"password"`

	// Version gate kept for config compatibility with the manual resource.
	PasswordWoVersion types.Int64 `tfsdk:"password_wo_version" apiname:"passwordWoVersion" autogen:"omitjson"`
}

var _ autogen.ResourceSchemaHook = (*rs)(nil)
var _ autogen.PostReadAPICallHook = (*rs)(nil)
var _ autogen.PostCreateAPICallHook = (*rs)(nil)

func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	if passwordAttr, ok := baseSchema.Attributes["password"].(schema.StringAttribute); ok {
		passwordAttr.Validators = append(passwordAttr.Validators, stringvalidator.ConflictsWith(path.Expressions{
			path.MatchRelative().AtParent().AtName("x509_type"),
			path.MatchRelative().AtParent().AtName("ldap_auth_type"),
			path.MatchRelative().AtParent().AtName("aws_iam_type"),
			path.MatchRelative().AtParent().AtName("password_wo"),
		}...))
		baseSchema.Attributes["password"] = passwordAttr
	}
	baseSchema.Attributes["password_wo"] = schema.StringAttribute{
		Optional:  true,
		Sensitive: true,
		WriteOnly: true,
		Validators: []validator.String{
			stringvalidator.ConflictsWith(path.Expressions{
				path.MatchRelative().AtParent().AtName("password"),
				path.MatchRelative().AtParent().AtName("x509_type"),
				path.MatchRelative().AtParent().AtName("ldap_auth_type"),
				path.MatchRelative().AtParent().AtName("aws_iam_type"),
			}...),
			stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("password_wo_version")),
		},
	}
	baseSchema.Attributes["password_wo_version"] = schema.Int64Attribute{
		Optional: true,
		Validators: []validator.Int64{
			int64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("password_wo")),
		},
	}
	for _, name := range []string{"project_id", "username"} {
		if attr, ok := baseSchema.Attributes[name].(schema.StringAttribute); ok {
			attr.PlanModifiers = append(attr.PlanModifiers, stringplanmodifier.RequiresReplace())
			baseSchema.Attributes[name] = attr
		}
	}
	// Crafted id: kept for migrated state (manual-era state carries the encoded id) and for
	// tooling that references .id; populated by the read/create hooks below.
	baseSchema.Attributes["id"] = schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	return baseSchema
}

// PostReadAPICall reproduces the manual resource's crafted id so state created by the
// manual implementation keeps resolving to the same identity (lift-and-shift).
func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}
	model, ok := req.State.(*TFModel)
	if !ok || model.ProjectId.IsNull() || model.Username.IsNull() || model.AuthDatabaseName.IsNull() {
		return result
	}

	var obj map[string]any
	if err := autogen.Decode(result.Body, &obj); err != nil {
		return result
	}
	obj["id"] = conversion.EncodeStateID(map[string]string{
		"project_id":         model.ProjectId.ValueString(),
		"username":           model.Username.ValueString(),
		"auth_database_name": model.AuthDatabaseName.ValueString(),
	})
	// Optional-string shape parity: the API returns "" for unset optional strings while
	// the manual resource keeps them null in state; drop them so state stays null
	// (inverse of the past migration's SDKv2 normalization, which set "" instead).
	if model.Description.IsNull() {
		if val, exists := obj["description"]; exists && val == "" {
			delete(obj, "description")
		}
	}
	body, err := json.Marshal(obj)
	if err != nil {
		return result
	}
	result.Body = body
	return result
}

// PostCreateAPICall injects the crafted id after create: the database users POST response
// carries no id field, but groupId/username/databaseName are present in the response body.
func (r *rs) PostCreateAPICall(result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil || len(result.Body) == 0 {
		return result
	}
	var obj map[string]any
	if err := autogen.Decode(result.Body, &obj); err != nil {
		return result
	}
	groupID, _ := obj["groupId"].(string)
	username, _ := obj["username"].(string)
	databaseName, _ := obj["databaseName"].(string)
	if groupID == "" || username == "" || databaseName == "" {
		return result
	}
	obj["id"] = conversion.EncodeStateID(map[string]string{
		"project_id":         groupID,
		"username":           username,
		"auth_database_name": databaseName,
	})
	body, err := json.Marshal(obj)
	if err != nil {
		return result
	}
	result.Body = body
	return result
}
