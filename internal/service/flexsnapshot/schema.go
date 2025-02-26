package flexsnapshot

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
			"snapshot_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the snapshot to restore.",
			},
			"expiration": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when the download link no longer works. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"finish_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud completed writing this snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"mongo_db_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "MongoDB host version that the snapshot runs.",
			},
			"scheduled_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud will take the snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"start_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud began taking the snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Phase of the restore workflow for this job at the time this resource made this request.",
			},
		},
	}
}

type TFModel struct {
	Expiration     types.String `tfsdk:"expiration"`
	FinishTime     types.String `tfsdk:"finish_time"`
	ProjectId      types.String `tfsdk:"project_id"`
	MongoDBVersion types.String `tfsdk:"mongo_db_version"`
	Name           types.String `tfsdk:"name"`
	ScheduledTime  types.String `tfsdk:"scheduled_time"`
	SnapshotId     types.String `tfsdk:"snapshot_id"`
	StartTime      types.String `tfsdk:"start_time"`
	Status         types.String `tfsdk:"status"`
}

type TFFlexSnapshotsDSModel struct {
	ProjectId types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	Results   []TFModel    `tfsdk:"results"`
}
