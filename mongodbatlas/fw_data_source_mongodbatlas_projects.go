package mongodbatlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20230201002/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ProjectsDS{}
var _ datasource.DataSourceWithConfigure = &ProjectsDS{}

var _ datasource.DataSource = &ProjectsDS{}
var _ datasource.DataSourceWithConfigure = &ProjectsDS{}

func NewProjectsDS() datasource.DataSource {
	return &ProjectsDS{}
}

type ProjectsDS struct {
	client *MongoDBClient
}

type tfProjectsDSModel struct {
	ID           types.String       `tfsdk:"id"`
	Results      []tfProjectDSModel `tfsdk:"results"`
	PageNum      types.Int64        `tfsdk:"page_num"`
	ItemsPerPage types.Int64        `tfsdk:"items_per_page"`
	TotalCount   types.Int64        `tfsdk:"total_count"`
}

func (d *ProjectsDS) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectsDS) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*MongoDBClient)

	if !ok {
		resp.Diagnostics.AddError(errorConfigureSummary, fmt.Sprintf(errorConfigure, req.ProviderData))
		return
	}
	d.client = client
}

func (d *ProjectsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// https://github.com/hashicorp/terraform-plugin-testing/issues/84#issuecomment-1480006432
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Please use each project's id attribute instead",
				MarkdownDescription: "Please use each project's id attribute instead",
				Computed:            true,
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
					},
				},
			},
		},
	}
}

func (d *ProjectsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var stateModel tfProjectsDSModel
	conn := d.client.Atlas
	connV2 := d.client.AtlasV2

	resp.Diagnostics.Append(req.Config.Get(ctx, &stateModel)...)
	options := &matlas.ListOptions{
		PageNum:      int(stateModel.PageNum.ValueInt64()),
		ItemsPerPage: int(stateModel.ItemsPerPage.ValueInt64()),
	}

	projectsRes, _, err := conn.Projects.GetAllProjects(ctx, options)
	if err != nil {
		resp.Diagnostics.AddError("error in monogbatlas_projects data source", fmt.Sprintf("error getting projects information: %s", err.Error()))
		return
	}

	err = populateProjectsDataSourceModel(ctx, conn, connV2, &stateModel, projectsRes)
	if err != nil {
		resp.Diagnostics.AddError("error in monogbatlas_projects data source", err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func populateProjectsDataSourceModel(ctx context.Context, conn *matlas.Client, connV2 *admin.APIClient, stateModel *tfProjectsDSModel, projectsRes *matlas.Projects) error {
	results := make([]tfProjectDSModel, len(projectsRes.Results))

	for i, project := range projectsRes.Results {
		atlasTeams, atlasLimits, atlasProjectSettings, err := getProjectPropsFromAPI(ctx, conn, connV2, project.ID)
		if err != nil {
			return fmt.Errorf("error while getting project properties for project %s: %v", project.ID, err.Error())
		}
		projectModel := toTFProjectDataSourceModel(ctx, project, atlasTeams, atlasProjectSettings, atlasLimits)
		results[i] = projectModel
	}

	stateModel.Results = results
	stateModel.TotalCount = types.Int64Value(int64(projectsRes.TotalCount))
	stateModel.ID = types.StringValue("test")
	return nil
}
