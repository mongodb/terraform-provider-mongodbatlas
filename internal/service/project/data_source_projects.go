package project

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312004/admin"
)

const projectsDataSourceName = "projects"

var _ datasource.DataSource = &ProjectsDS{}
var _ datasource.DataSourceWithConfigure = &ProjectsDS{}

func PluralDataSource() datasource.DataSource {
	return &ProjectsDS{
		DSCommon: config.DSCommon{
			DataSourceName: projectsDataSourceName,
		},
	}
}

type ProjectsDS struct {
	config.DSCommon
}

type tfProjectsDSModel struct {
	ID           types.String        `tfsdk:"id"`
	Results      []*TFProjectDSModel `tfsdk:"results"`
	PageNum      types.Int64         `tfsdk:"page_num"`
	ItemsPerPage types.Int64         `tfsdk:"items_per_page"`
	TotalCount   types.Int64         `tfsdk:"total_count"`
}

func (d *ProjectsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// https://github.com/hashicorp/terraform-plugin-testing/issues/84#issuecomment-1480006432
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage: "Please use each project's id attribute instead",
				Computed:           true,
			},
			"page_num": schema.Int64Attribute{
				Optional: true,
			},
			"items_per_page": schema.Int64Attribute{
				Optional: true,
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"org_id": schema.StringAttribute{
							Computed: true,
						},
						"project_id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"cluster_count": schema.Int64Attribute{
							Computed: true,
						},
						"created": schema.StringAttribute{
							Computed: true,
						},
						"is_collect_database_specifics_statistics_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"is_data_explorer_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"is_extended_storage_sizes_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"is_performance_advisor_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"is_realtime_performance_panel_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"is_schema_advisor_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"is_slow_operation_thresholding_enabled": schema.BoolAttribute{
							DeprecationMessage: constant.DeprecationParam, // added deprecation in CLOUDP-293855 because was deprecated in the doc
							Computed:           true,
						},
						"region_usage_restrictions": schema.StringAttribute{
							Computed: true,
						},
						"teams": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"team_id": schema.StringAttribute{
										Computed: true,
									},
									"role_names": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
						"limits": schema.SetNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed: true,
									},
									"value": schema.Int64Attribute{
										Computed: true,
									},
									"current_usage": schema.Int64Attribute{
										Computed: true,
									},
									"default_limit": schema.Int64Attribute{
										Computed: true,
									},
									"maximum_limit": schema.Int64Attribute{
										Computed: true,
									},
								},
							},
						},
						"ip_addresses": schema.SingleNestedAttribute{
							Computed:           true,
							DeprecationMessage: fmt.Sprintf(constant.DeprecationParamWithReplacement, "mongodbatlas_project_ip_addresses data source"),
							Attributes: map[string]schema.Attribute{
								"services": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"clusters": schema.ListNestedAttribute{
											Computed: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"cluster_name": schema.StringAttribute{
														Computed: true,
													},
													"inbound": schema.ListAttribute{
														ElementType: types.StringType,
														Computed:    true,
													},
													"outbound": schema.ListAttribute{
														ElementType: types.StringType,
														Computed:    true,
													},
												},
											},
										},
									},
								},
							},
						},
						"tags": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"users": schema.SetNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"org_membership_status": schema.StringAttribute{
										Computed: true,
									},
									"roles": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
									"username": schema.StringAttribute{
										Computed: true,
									},
									"invitation_created_at": schema.StringAttribute{
										Computed: true,
									},
									"invitation_expires_at": schema.StringAttribute{
										Computed: true,
									},
									"inviter_username": schema.StringAttribute{
										Computed: true,
									},
									"country": schema.StringAttribute{
										Computed: true,
									},
									"created_at": schema.StringAttribute{
										Computed: true,
									},
									"first_name": schema.StringAttribute{
										Computed: true,
									},
									"last_auth": schema.StringAttribute{
										Computed: true,
									},
									"last_name": schema.StringAttribute{
										Computed: true,
									},
									"mobile_number": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *ProjectsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var stateModel tfProjectsDSModel
	connV2 := d.Client.AtlasV2

	resp.Diagnostics.Append(req.Config.Get(ctx, &stateModel)...)

	projectParams := &admin.ListProjectsApiParams{
		PageNum:      conversion.IntPtr(int(stateModel.PageNum.ValueInt64())),
		ItemsPerPage: conversion.IntPtr(int(stateModel.ItemsPerPage.ValueInt64())),
	}
	projectsRes, _, err := connV2.ProjectsApi.ListProjectsWithParams(ctx, projectParams).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error in monogbatlas_projects data source", fmt.Sprintf("error getting projects information: %s", err.Error()))
		return
	}

	diags := populateProjectsDataSourceModel(ctx, connV2, &stateModel, projectsRes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func populateProjectsDataSourceModel(ctx context.Context, connV2 *admin.APIClient, stateModel *tfProjectsDSModel, projectsRes *admin.PaginatedAtlasGroup) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	input := projectsRes.GetResults()
	results := make([]*TFProjectDSModel, 0, len(input))
	for i := range input {
		project := input[i]

		projectPropsParams := &PropsParams{
			ProjectID:             project.GetId(),
			IsDataSource:          true,
			ProjectsAPI:           connV2.ProjectsApi,
			TeamsAPI:              connV2.TeamsApi,
			PerformanceAdvisorAPI: connV2.PerformanceAdvisorApi,
			MongoDBCloudUsersAPI:  connV2.MongoDBCloudUsersApi,
		}

		projectProps, err := GetProjectPropsFromAPI(ctx, projectPropsParams, &diagnostics)
		if err == nil { // if the project is still valid, e.g. could have just been deleted
			projectModel, diags := NewTFProjectDataSourceModel(ctx, &project, projectProps)
			diagnostics = append(diagnostics, diags...)
			if projectModel != nil {
				results = append(results, projectModel)
			}
		}
	}
	stateModel.Results = results
	stateModel.TotalCount = types.Int64Value(int64(projectsRes.GetTotalCount()))
	stateModel.ID = types.StringValue(id.UniqueId())
	return diagnostics
}
