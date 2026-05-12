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

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dynamicjson"
)

func UpdateResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Patch-only generic Terraform resource. Sets fields on an Atlas entity that " +
			"is primarily managed by another (typed) resource — without taking over its lifecycle. " +
			"Use this for preview fields the typed resource does not yet expose. " +
			"WARNING: the `output` attribute contains the full API response and is not marked Sensitive. " +
			"Any secret returned by the API will appear in plan/apply output unless piped through a sensitive Terraform `output` block.",
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
			"output": schema.DynamicAttribute{
				Computed: true,
				Description: "Full API response from the most recent PATCH or GET. Access fields with " +
					"`mongodbatlas_api_update.<name>.output.<field>`. NOT marked Sensitive.",
			},
		},
	}
}
