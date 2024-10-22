package flexcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: dataSourceSchema(false),
	}
}

func dataSourceSchema(isPlural bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Required:            !isPlural,
			Computed:            isPlural,
			MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
		},
		"name": schema.StringAttribute{
			Required:            !isPlural,
			Computed:            isPlural,
			MarkdownDescription: "Human-readable label that identifies the flex cluster.",
		},
		"provider_settings": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"backing_provider_name": schema.StringAttribute{
					Computed:            true,
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
					Computed:            true,
					MarkdownDescription: "Human-readable label that identifies the geographic location of your MongoDB flex cluster. The region you choose can affect network latency for clients accessing your databases. For a complete list of region names, see [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/#std-label-amazon-aws), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), and [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).",
				},
			},
			Computed:            true,
			MarkdownDescription: "Group of cloud provider settings that configure the provisioned MongoDB flex cluster.",
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
		"mongo_dbversion": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Version of MongoDB that the instance runs.",
		},
		"state_name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Human-readable label that indicates the current operating condition of this instance.",
		},
		"tags": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Constant that defines the set of the tag. For example, `environment` in the `environment : production` tag.",
					},
					"value": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Variable that belongs to the set of the tag. For example, `production` in the `environment : production` tag.",
					},
				},
			},
			Computed:            true,
			MarkdownDescription: "List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the instance.",
		},
		"termination_protection_enabled": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.",
		},
		"version_release_system": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Method by which the cluster maintains the MongoDB versions.",
		},
	}
}
