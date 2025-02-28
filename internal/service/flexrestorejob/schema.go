package flexrestorejob

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable label that identifies the flex cluster whose snapshot you want to restore.",
			},
			"restore_job_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the restore job.",
			},
			"delivery_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Means by which this resource returns the snapshot to the requesting MongoDB Cloud user.",
			},
			"expiration_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when the download link no longer works. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"restore_finished_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud completed writing this snapshot. MongoDB Cloud changes the status of the restore job to `CLOSED`. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"restore_scheduled_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud will restore this snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"snapshot_finished_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud completed writing this snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"snapshot_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the snapshot to restore.",
			},
			"snapshot_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internet address from which you can download the compressed snapshot files. The resource returns this parameter when  `\"deliveryType\" : \"DOWNLOAD\"`.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Phase of the restore workflow for this job at the time this resource made this request.",
			},
			"target_deployment_item_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Human-readable label that identifies the instance or cluster on the target project to which you want to restore the snapshot. You can restore the snapshot to another flex cluster or dedicated cluster tier.",
			},
			"target_project_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the project that contains the instance or cluster to which you want to restore the snapshot.",
			},
		},
	}
}

type TFModel struct {
	DeliveryType             types.String `tfsdk:"delivery_type"`
	ExpirationDate           types.String `tfsdk:"expiration_date"`
	ProjectID                types.String `tfsdk:"project_id"`
	RestoreJobID             types.String `tfsdk:"restore_job_id"`
	Name                     types.String `tfsdk:"name"`
	RestoreFinishedDate      types.String `tfsdk:"restore_finished_date"`
	RestoreScheduledDate     types.String `tfsdk:"restore_scheduled_date"`
	SnapshotFinishedDate     types.String `tfsdk:"snapshot_finished_date"`
	SnapshotID               types.String `tfsdk:"snapshot_id"`
	SnapshotUrl              types.String `tfsdk:"snapshot_url"`
	Status                   types.String `tfsdk:"status"`
	TargetDeploymentItemName types.String `tfsdk:"target_deployment_item_name"`
	TargetProjectID          types.String `tfsdk:"target_project_id"`
}

type TFFlexRestoreJobsDSModel struct {
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	Results   []TFModel    `tfsdk:"results"`
}
