package apiresource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dynamicjson"
)

func ResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Generic Terraform resource wrapping any Atlas Admin API endpoint. " +
			"Useful for endpoints not yet covered by a typed resource. " +
			"WARNING: the `output` attribute contains the full API response and is not marked Sensitive. " +
			"Any secret returned by the API will appear in plan/apply output unless piped through a sensitive Terraform `output` block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Synthetic identifier set to the resolved read URL.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Required: true,
				Description: "Collection path used for Create (e.g. `/api/atlas/v2/orgs/<orgId>/serviceAccounts`). " +
					"For singleton endpoints, also used as-is for Read/Update/Delete. " +
					"For collection endpoints, the provider appends `/<id>` from `output` (see `id_attribute`) for Read/Update/Delete.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id_attribute": schema.ListAttribute{
				Optional:    true,
				ElementType: basetypes.StringType{},
				Description: "Response field name(s) that identify the resource. The provider reads each name from `output` " +
					"and appends `/<value>` to `path` for Read/Update/Delete (in order, for composite IDs). " +
					"Omit for singleton endpoints where Create/Read/Update/Delete share the same path. " +
					"Example: `[\"clientId\"]` for service accounts; `[\"databaseName\", \"username\"]` for database users.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"create_method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(defaultCreateMethod),
				Description: "HTTP method for Create. One of POST, PUT, PATCH. Defaults to POST.",
				Validators:  []validator.String{stringvalidator.OneOf("POST", "PUT", "PATCH")},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"update_method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(defaultUpdateMethod),
				Description: "HTTP method for Update. One of PATCH, PUT, POST. Defaults to PATCH.",
				Validators:  []validator.String{stringvalidator.OneOf("PATCH", "PUT", "POST")},
			},
			"version_header": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Atlas API version media type used for both Accept and Content-Type headers. " +
					"When unset, defaults to today's UTC date in the form `application/vnd.atlas.<YYYY-MM-DD>+json` " +
					"(Atlas snaps the date down to the latest published version on or before it). The resolved value " +
					"is persisted at Create time and stays stable for the lifetime of the resource. " +
					"Mutually exclusive with `preview`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"preview": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Shorthand for `version_header = \"" + previewVersionHeader + "\"`. Mutually exclusive with `version_header`.",
			},
			"body": schema.DynamicAttribute{
				Optional: true,
				Computed: true,
				Description: "Request body sent on Create/Update. The reshape engine drives drift detection: " +
					"keys you declare here are tracked; fields Atlas adds server-side are stored in `output`.",
				PlanModifiers: []planmodifier.Dynamic{
					customplanmodifier.DynamicUseStateWhen(dynamicjson.SemanticallyEqual),
				},
			},
			"sensitive_body": schema.DynamicAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "Sensitive fragment merged into the request body at Create/Update. " +
					"Values are never written to state; their keys are excluded from drift comparison.",
			},
			"create_only_body_keys": schema.SetAttribute{
				Optional:    true,
				ElementType: basetypes.StringType{},
				Description: "Top-level body keys that the endpoint accepts only on Create and rejects on Update " +
					"(for example `secretExpiresAfterHours`). These keys are stripped from the payload " +
					"before the Update request is issued.",
			},
			"output": schema.DynamicAttribute{
				Computed: true,
				Description: "Full API response from the most recent operation. Access fields with " +
					"`mongodbatlas_api_resource.<name>.output.<field>`. NOT marked Sensitive — see the resource-level warning.",
			},
		},
	}
}
