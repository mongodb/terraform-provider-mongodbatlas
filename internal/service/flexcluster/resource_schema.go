package flexcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal character string that identifies the project.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable label that identifies the instance.",
			},
			"provider_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"backing_provider_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Cloud service provider on which MongoDB Cloud provisioned the flex cluster.",
					},
					"disk_size_gb": schema.Float64Attribute{
						Computed:            true,
						MarkdownDescription: "Storage capacity available to the flex cluster expressed in gigabytes.",
					},
					"provider_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Human-readable label that identifies the cloud service provider.",
					},
					"region_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Human-readable label that identifies the geographic location of your MongoDB flex cluster. The region you choose can affect network latency for clients accessing your databases. For a complete list of region names, see [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/#std-label-amazon-aws), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), and [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).",
					},
				},
				Required:            true,
				MarkdownDescription: "Group of cloud provider settings that configure the provisioned MongoDB flex cluster.",
			},
			"tags": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Map that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the instance.",
			},
			"backup_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Flag that indicates whether backups are performed for this flex cluster. Backup uses [TODO](TODO) for flex clusters.",
					},
				},
				Computed:            true,
				MarkdownDescription: "Flex backup configuration",
			},
			"cluster_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Flex cluster topology.",
			},
			"connection_strings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"standard": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Public connection string that you can use to connect to this cluster. This connection string uses the mongodb:// protocol.",
					},
					"standard_srv": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Public connection string that you can use to connect to this flex cluster. This connection string uses the `mongodb+srv://` protocol.",
					},
				},
				Computed:            true,
				MarkdownDescription: "Collection of Uniform Resource Locators that point to the MongoDB database.",
			},
			"create_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud created this instance. This parameter expresses its value in ISO 8601 format in UTC.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the instance.",
			},
			"mongo_db_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Version of MongoDB that the instance runs.",
			},
			"state_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Human-readable label that indicates the current operating condition of this instance.",
			},
			"termination_protection_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.",
			},
			"version_release_system": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Method by which the cluster maintains the MongoDB versions.",
			},
		},
	}
}

type TFModel struct {
	ProviderSettings             types.Object `tfsdk:"provider_settings"`
	ConnectionStrings            types.Object `tfsdk:"connection_strings"`
	Tags                         types.Map    `tfsdk:"tags"`
	CreateDate                   types.String `tfsdk:"create_date"`
	ProjectId                    types.String `tfsdk:"project_id"`
	Id                           types.String `tfsdk:"id"`
	MongoDbversion               types.String `tfsdk:"mongo_db_version"`
	Name                         types.String `tfsdk:"name"`
	ClusterType                  types.String `tfsdk:"cluster_type"`
	StateName                    types.String `tfsdk:"state_name"`
	VersionReleaseSystem         types.String `tfsdk:"version_release_system"`
	BackupSettings               types.Object `tfsdk:"backup_settings"`
	TerminationProtectionEnabled types.Bool   `tfsdk:"termination_protection_enabled"`
}

type TFBackupSettings struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

var BackupSettingsType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"enabled": types.BoolType,
}}

type TFConnectionStrings struct {
	Standard    types.String `tfsdk:"standard"`
	StandardSrv types.String `tfsdk:"standard_srv"`
}

var ConnectionStringsType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"standard":     types.StringType,
	"standard_srv": types.StringType,
}}

type TFProviderSettings struct {
	BackingProviderName types.String  `tfsdk:"backing_provider_name"`
	DiskSizeGb          types.Float64 `tfsdk:"disk_size_gb"`
	ProviderName        types.String  `tfsdk:"provider_name"`
	RegionName          types.String  `tfsdk:"region_name"`
}

var ProviderSettingsType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"backing_provider_name": types.StringType,
	"disk_size_gb":          types.Float64Type,
	"provider_name":         types.StringType,
	"region_name":           types.StringType,
}}
