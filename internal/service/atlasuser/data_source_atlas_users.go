package atlasuser

import (
	"context"
	"fmt"

	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	AtlasUsersDataSourceName      = "atlas_users"
	errorUsersRead                = "error getting atlas users(%s - %s): %s"
	errorMissingAttributesSummary = "missing required attributes for data source"
	ErrorMissingAttributesDetail  = "either org_id, project_id, or team_id with org_id must be configured"
)

var _ datasource.DataSource = &atlasUsersDS{}
var _ datasource.DataSourceWithConfigure = &atlasUsersDS{}

func PluralDataSource() datasource.DataSource {
	return &atlasUsersDS{
		DSCommon: config.DSCommon{
			DataSourceName: AtlasUsersDataSourceName,
		},
	}
}

type atlasUsersDS struct {
	config.DSCommon
}

type tfAtlasUsersDSModel struct {
	ID           types.String         `tfsdk:"id"`
	OrgID        types.String         `tfsdk:"org_id"`
	ProjectID    types.String         `tfsdk:"project_id"`
	TeamID       types.String         `tfsdk:"team_id"`
	Results      []tfAtlasUserDSModel `tfsdk:"results"`
	PageNum      types.Int64          `tfsdk:"page_num"`
	ItemsPerPage types.Int64          `tfsdk:"items_per_page"`
	TotalCount   types.Int64          `tfsdk:"total_count"`
}

func (d *atlasUsersDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		DeprecationMessage: fmt.Sprintf(constant.DeprecationNextMajorWithReplacementGuide, "data source", "data.mongodbatlas_organization.users, data.mongodbatlas_team.users or data.mongodbatlas_project.users attributes", "[Migration Guide: Migrate off deprecated `mongodbatlas_atlas_user` and `mongodbatlas_atlas_users`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/atlas-user-migration-guide)"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework: https://github.com/hashicorp/terraform-plugin-testing/issues/84#issuecomment-1480006432
				DeprecationMessage: "Please use each user's id attribute instead",
				Computed:           true,
			},
			"org_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("project_id")),
				},
			},
			"project_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("org_id"),
						path.MatchRoot("team_id"),
					}...),
				},
			},
			"team_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("project_id")),
				},
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
						"user_id": schema.StringAttribute{
							Computed: true,
						},
						"username": schema.StringAttribute{
							Computed: true,
						},
						"country": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"email_address": schema.StringAttribute{
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
						"team_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"links": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"href": schema.StringAttribute{
										Computed: true,
									},
									"rel": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"roles": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"group_id": schema.StringAttribute{
										Computed: true,
									},
									"org_id": schema.StringAttribute{
										Computed: true,
									},
									"role_name": schema.StringAttribute{
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

func (d *atlasUsersDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	connV2 := d.Client.AtlasV220241113

	var atlasUsersConfig tfAtlasUsersDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &atlasUsersConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if atlasUsersConfig.OrgID.IsNull() && atlasUsersConfig.ProjectID.IsNull() {
		resp.Diagnostics.AddError(errorMissingAttributesSummary, ErrorMissingAttributesDetail)
		return
	}

	var (
		users      []admin20241113.CloudAppUser
		totalCount int
	)

	switch {
	case !atlasUsersConfig.ProjectID.IsNull():
		projectID := atlasUsersConfig.ProjectID.ValueString()
		apiResp, _, err := connV2.ProjectsApi.ListProjectUsersWithParams(ctx, &admin20241113.ListProjectUsersApiParams{
			GroupId:      projectID,
			PageNum:      conversion.Int64PtrToIntPtr(atlasUsersConfig.PageNum.ValueInt64Pointer()),
			ItemsPerPage: conversion.Int64PtrToIntPtr(atlasUsersConfig.ItemsPerPage.ValueInt64Pointer()),
		}).Execute()
		if err != nil {
			resp.Diagnostics.AddError("error when getting users from Atlas", fmt.Sprintf(errorUsersRead, "project", projectID, err.Error()))
			return
		}
		users = apiResp.GetResults()
		totalCount = *apiResp.TotalCount
	case !atlasUsersConfig.TeamID.IsNull() && !atlasUsersConfig.OrgID.IsNull():
		teamID := atlasUsersConfig.TeamID.ValueString()
		apiResp, _, err := connV2.TeamsApi.ListTeamUsersWithParams(ctx, &admin20241113.ListTeamUsersApiParams{
			OrgId:        atlasUsersConfig.OrgID.ValueString(),
			TeamId:       teamID,
			PageNum:      conversion.Int64PtrToIntPtr(atlasUsersConfig.PageNum.ValueInt64Pointer()),
			ItemsPerPage: conversion.Int64PtrToIntPtr(atlasUsersConfig.ItemsPerPage.ValueInt64Pointer()),
		}).Execute()
		if err != nil {
			resp.Diagnostics.AddError("error when getting users from Atlas", fmt.Sprintf(errorUsersRead, "team", teamID, err.Error()))
			return
		}
		users = apiResp.GetResults()
		totalCount = *apiResp.TotalCount
	default: // only org_id is defined
		orgID := atlasUsersConfig.OrgID.ValueString()
		apiResp, _, err := connV2.OrganizationsApi.ListOrganizationUsersWithParams(ctx, &admin20241113.ListOrganizationUsersApiParams{
			OrgId:        atlasUsersConfig.OrgID.ValueString(),
			PageNum:      conversion.Int64PtrToIntPtr(atlasUsersConfig.PageNum.ValueInt64Pointer()),
			ItemsPerPage: conversion.Int64PtrToIntPtr(atlasUsersConfig.ItemsPerPage.ValueInt64Pointer()),
		}).Execute()
		if err != nil {
			resp.Diagnostics.AddError("error when getting users from Atlas", fmt.Sprintf(errorUsersRead, "org", orgID, err.Error()))
			return
		}
		users = apiResp.GetResults()
		totalCount = *apiResp.TotalCount
	}

	usersResultState := newTFAtlasUsersDSModel(&atlasUsersConfig, users, totalCount)
	resp.Diagnostics.Append(resp.State.Set(ctx, &usersResultState)...)
}

func newTFAtlasUsersDSModel(atlasUsersConfig *tfAtlasUsersDSModel, users []admin20241113.CloudAppUser, totalCount int) tfAtlasUsersDSModel {
	return tfAtlasUsersDSModel{
		ID:           types.StringValue(id.UniqueId()),
		OrgID:        atlasUsersConfig.OrgID,
		ProjectID:    atlasUsersConfig.ProjectID,
		TeamID:       atlasUsersConfig.TeamID,
		PageNum:      atlasUsersConfig.PageNum,
		ItemsPerPage: atlasUsersConfig.ItemsPerPage,
		TotalCount:   types.Int64Value(int64(totalCount)),
		Results:      newTFAtlasUsersList(users),
	}
}

func newTFAtlasUsersList(users []admin20241113.CloudAppUser) []tfAtlasUserDSModel {
	resUsers := make([]tfAtlasUserDSModel, len(users))
	for i := range users {
		resUsers[i] = newTFAtlasUserDSModel(&users[i])
	}
	return resUsers
}
