package organization3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "EXPERIMENTAL PoC resource demonstrating ModifyPlan-driven client secret rotation against the Atlas API. " +
			"Do not use in production; see POC_README.md.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Organization display name. Also used as the service account name at create.",
				Required:            true,
			},
			"org_owner_id": schema.StringAttribute{
				MarkdownDescription: "Atlas user ID of the organization owner. Required at create.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Atlas organization ID assigned at create.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Service account client ID from organization create.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Service account client secret. Sensitive; set at create only and not refreshed on read.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_secret_rotation": schema.SingleNestedAttribute{
				MarkdownDescription: "When set, opts into ModifyPlan-driven secret rotation using Atlas `expires_at` metadata.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"expires_after_hours": schema.Int64Attribute{
						MarkdownDescription: "Secret lifetime in hours for the next secret created. Applies on create and each rotation POST; does not change existing secrets.",
						Required:            true,
					},
					"rotate_before_expiry_hours": schema.Int64Attribute{
						MarkdownDescription: "Hours before `current_secret.expires_at` when ModifyPlan should schedule rotation. Defaults to `expires_after_hours / 2`.",
						Optional:            true,
						Computed:            true,
					},
					"secret_version": schema.Int64Attribute{
						MarkdownDescription: "Rotation generation. ModifyPlan may increment when renewal is due; set higher to force rotation.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"current_secret": schema.SingleNestedAttribute{
						MarkdownDescription: "Active secret metadata refreshed from Atlas on read.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"secret_id": schema.StringAttribute{
								Computed: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"created_at": schema.StringAttribute{
								MarkdownDescription: "RFC3339 timestamp from Atlas.",
								Computed:            true,
							},
							"expires_at": schema.StringAttribute{
								MarkdownDescription: "RFC3339 expiry from Atlas; source of truth for rotation scheduling.",
								Computed:            true,
							},
							"last_used_at": schema.StringAttribute{
								MarkdownDescription: "RFC3339 last use from Atlas. Null when unused.",
								Computed:            true,
							},
						},
					},
					"old_secret": schema.SingleNestedAttribute{
						MarkdownDescription: "Previous secret metadata after rotation. Null when no overlap secret exists.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"secret_id": schema.StringAttribute{
								Computed: true,
							},
							"created_at": schema.StringAttribute{
								Computed: true,
							},
							"expires_at": schema.StringAttribute{
								Computed: true,
							},
							"last_used_at": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}
