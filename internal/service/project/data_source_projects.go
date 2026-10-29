package project

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
)

const projectsDataSourceName = "projects"

var _ datasource.DataSource = &ProjectsDS{}
var _ datasource.DataSourceWithConfigure = &ProjectsDS{}

func PluralDataSource() datasource.DataSource {
	return &ProjectsDS{
		DataSourceName: projectsDataSourceName,
	}
}

type ProjectsDS struct {
	config.DSCommon
}

func (d *ProjectsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(), &conversion.PluralDataSourceSchemaRequest{
		OverridenFields: dataSourceOverridenFields(),
		HasLegacyFields: true,
	})
}

func (d *ProjectsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var stateModel tfProjectsDSModel
	connV2 := d.Client.AtlasV2

	resp.Diagnostics.Append(req.Config.Get(ctx, &stateModel)...)

	projectParams := &admin.ListGroupsApiParams{
		PageNum:      conversion.IntPtr(int(stateModel.PageNum.ValueInt64())),
		ItemsPerPage: conversion.IntPtr(int(stateModel.ItemsPerPage.ValueInt64())),
	}
	projectsRes, _, err := connV2.ProjectsAPI.ListGroupsWithParams(ctx, projectParams).Execute()
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
			ProjectsAPI:           connV2.ProjectsAPI,
			TeamsAPI:              connV2.TeamsAPI,
			PerformanceAdvisorAPI: connV2.PerformanceAdvisorAPI,
			MongoDBCloudUsersAPI:  connV2.MongoDBCloudUsersAPI,
		}

		projectProps, err := GetProjectPropsFromAPI(ctx, projectPropsParams)
		if err != nil {
			if isNotFoundError(err) { // the project could have been deleted after it was listed
				continue
			}
			diagnostics.AddError(fmt.Sprintf("error getting properties of project %s", project.GetId()), err.Error())
			return diagnostics
		}
		projectModel, diags := NewTFProjectDataSourceModel(ctx, &project, projectProps)
		diagnostics = append(diagnostics, diags...)
		if projectModel != nil {
			results = append(results, projectModel)
		}
	}
	stateModel.Results = results
	stateModel.TotalCount = types.Int64Value(int64(projectsRes.GetTotalCount()))
	stateModel.ID = types.StringValue(id.UniqueId())
	return diagnostics
}

// isNotFoundError reports whether err is an Atlas API 404 response.
func isNotFoundError(err error) bool {
	apiErr, ok := admin.AsError(err)
	return ok && apiErr.GetError() == http.StatusNotFound
}
