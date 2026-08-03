package project

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312022/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &projectDS{}
var _ datasource.DataSourceWithConfigure = &projectDS{}

func DataSource() datasource.DataSource {
	return &projectDS{
		DSCommon: config.DSCommon{
			DataSourceName: projectResourceName,
		},
	}
}

type projectDS struct {
	config.DSCommon
}

func (d *projectDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(), &conversion.DataSourceSchemaRequest{
		OverridenFields: dataSourceOverridenFields(),
	})
}

func (d *projectDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var projectState TFProjectDSModel
	connV2 := d.Client.AtlasV2

	resp.Diagnostics.Append(req.Config.Get(ctx, &projectState)...)

	if projectState.ProjectID.IsNull() && projectState.Name.IsNull() {
		resp.Diagnostics.AddError("missing required attributes for data source", "either project_id or name must be configured")
		return
	}

	var (
		err     error
		project *admin.Group
	)

	if !projectState.ProjectID.IsNull() {
		projectID := projectState.ProjectID.ValueString()
		project, _, err = connV2.ProjectsAPI.GetGroup(ctx, projectID).Execute()
		if err != nil {
			resp.Diagnostics.AddError("error when getting project from Atlas", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
			return
		}
	} else {
		name := projectState.Name.ValueString()
		project, _, err = connV2.ProjectsAPI.GetGroupByName(ctx, name).Execute()
		if err != nil {
			resp.Diagnostics.AddError("error when getting project from Atlas", fmt.Sprintf(ErrorProjectRead, name, err.Error()))
			return
		}
	}
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
		resp.Diagnostics.AddError("error when getting project properties", fmt.Sprintf(ErrorProjectRead, project.GetId(), err.Error()))
		return
	}

	newProjectState, diags := NewTFProjectDataSourceModel(ctx, project, projectProps)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newProjectState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ListAllProjectUsers(ctx context.Context, projectID string, mongoDBCloudUsersAPI admin.MongoDBCloudUsersAPI) ([]admin.GroupUserResponse, error) {
	return dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.GroupUserResponse], *http.Response, error) {
		request := mongoDBCloudUsersAPI.ListGroupUsers(ctx, projectID)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
}
