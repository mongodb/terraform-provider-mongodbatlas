package apiresource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dynamicjson"
)

func UpdateResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Patch-only generic Terraform resource. Sets fields on an Atlas entity that " +
			"is primarily managed by another (typed) resource — without taking over its lifecycle. " +
			"Use this for preview fields the typed resource does not yet expose. " +
			"By default no response fields are persisted in state — declare paths in `response_export_values` " +
			"(visible) or `response_export_values_sensitive` (redacted from plan/apply output) to opt in. " +
			"**Import is best-effort**: `terraform import` recovers only the resource URL; `body` and `sensitive_body` must be re-declared in HCL, and `sensitive_body` values cannot be recovered from Atlas (rotate or re-supply).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Synthetic identifier set to the resolved read URL (equal to `path`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Required: true,
				Description: "Full URL of the entity being patched (e.g. `/api/atlas/v2/groups/<gid>/streams/<name>`). " +
					"Used as-is for both PATCH (update) and GET (refresh). The entity must already exist.",
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
					"When unset, defaults to today's UTC date in the form `application/vnd.atlas.<YYYY-MM-DD>+json`. " +
					"The resolved value is persisted at first apply and stays stable for the lifetime of the resource. " +
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
				Description: "Request body sent on every Update. Only the keys you declare here are tracked for drift; " +
					"all other fields on the entity belong to the typed resource (or to Atlas server-side).",
				PlanModifiers: []planmodifier.Dynamic{
					customplanmodifier.DynamicUseStateWhen(dynamicjson.SemanticallyEqual),
				},
			},
			"sensitive_body": schema.DynamicAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "Sensitive fragment merged into the request body. Values are never written to state; " +
					"their keys are excluded from drift comparison.",
			},
			"response_export_values": schema.ListAttribute{
				Optional:    true,
				ElementType: basetypes.StringType{},
				Description: "Dotted paths into the API response to retain in `output`. Anything not listed " +
					"is discarded before state write. Numeric segments index lists (e.g. `failoverRegions.0.region`). " +
					"Missing paths are silently skipped.",
			},
			"response_export_values_sensitive": schema.ListAttribute{
				Optional:    true,
				ElementType: basetypes.StringType{},
				Description: "Same syntax as `response_export_values`, but matched values are stored in " +
					"`output_sensitive` (Sensitive). A path must not appear in both lists.",
			},
			"output": schema.DynamicAttribute{
				Computed: true,
				Description: "Projected API response. Contains only the paths listed in `response_export_values`. " +
					"Null when no paths are declared.",
			},
			"output_sensitive": schema.DynamicAttribute{
				Computed:  true,
				Sensitive: true,
				Description: "Projected API response containing the paths listed in `response_export_values_sensitive`. " +
					"Marked Sensitive: Terraform redacts values from plan/apply output. " +
					"Null when no sensitive paths are declared.",
			},
		},
	}
}
