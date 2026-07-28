package organization2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "EXPERIMENTAL PoC resource demonstrating ModifyPlan-driven client secret rotation. " +
			"Do not use in production; see POC_README.md.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable label for the mock organization.",
				Required:            true,
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Mock organization identifier assigned at create.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Mock service account client identifier.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Mock client secret. Sensitive; retained in state after create.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_secret_rotation": schema.SingleNestedAttribute{
				MarkdownDescription: "EXPERIMENTAL PoC nested object. When set, opts into ModifyPlan-driven secret rotation.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"interval": schema.StringAttribute{
						MarkdownDescription: "Duration between scheduled renewals (for example `2s`, `240h`).",
						Required:            true,
					},
					"secret_version": schema.Int64Attribute{
						MarkdownDescription: "Rotation generation. ModifyPlan may increment this when renewal is due; set a higher value to force rotation.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"next_renewal": schema.StringAttribute{
						MarkdownDescription: "RFC3339 timestamp when the next renewal is due.",
						Computed:            true,
					},
					"expires_at": schema.StringAttribute{
						MarkdownDescription: "RFC3339 timestamp when the current secret expires (2× interval after creation).",
						Computed:            true,
					},
					"current_secret_id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the active secret.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"old_secret_id": schema.StringAttribute{
						MarkdownDescription: "Previous secret identifier after rotation (teaching-only overlap visibility).",
						Computed:            true,
					},
				},
			},
		},
	}
}
