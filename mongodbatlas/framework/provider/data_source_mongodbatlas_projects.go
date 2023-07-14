package provider

import (
	"context"
	"fmt"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &ProjectsDataSource{}
var _ datasource.DataSourceWithConfigure = &ProjectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

type ProjectsDataSource struct {
	client *MongoDBClient
}

type projectsDataSourceModel struct {
	PageNum      types.Int64              `tfsdk:"page_num"`
	ItemsPerPage types.Int64              `tfsdk:"items_per_page"`
	Results      []projectDataSourceModel `tfsdk:"results"`
}

func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*MongoDBClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *MongoDBClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Info(ctx, "Schema() of example resource")

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"page_num": schema.Int64Attribute{
				Optional: true,
			},
			"items_per_page": schema.Int64Attribute{
				Optional: true,
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"project_id": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("name")),
							},
						},
						"name": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("project_id")),
							},
						},
						"org_id": schema.StringAttribute{
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
						"api_keys": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"api_key_id": schema.StringAttribute{
										Computed: true,
									},
									"role_names": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
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
					},
				},
			},
		},
	}
}

func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var stateModel projectsDataSourceModel
	conn := d.client.Atlas

	resp.Diagnostics.Append(req.Config.Get(ctx, &stateModel)...)
	options := &matlas.ListOptions{
		PageNum:      int(stateModel.PageNum.ValueInt64()),
		ItemsPerPage: int(stateModel.ItemsPerPage.ValueInt64()),
	}

	projects, _, err := conn.Projects.GetAllProjects(ctx, options)
	if err != nil {
		resp.Diagnostics.AddError("error in monogbatlas_projects data source", fmt.Sprintf("error getting projects information: %s", err.Error()))
		return
	}

	err = populateProjectsDataSourceModel(ctx, conn, stateModel, projects.Results)
	if err != nil {
		resp.Diagnostics.AddError("error in monogbatlas_projects data source", err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func populateProjectsDataSourceModel(ctx context.Context, conn *matlas.Client, stateModel projectsDataSourceModel, projects []*matlas.Project) error {
	results := make([]projectDataSourceModel, len(projects))

	for _, project := range projects {
		teams, apiKeys, projectSettings, err := getProjectPropsFromAtlas(ctx, conn, project)
		if err != nil {
			return fmt.Errorf("error while getting project properties for project %s: %v", project.ID, err.Error())
		}
		projectModel := toProjectDataSourceModel(ctx, project, teams, apiKeys, projectSettings)
		results = append(results, projectModel)
	}

	stateModel.Results = results
	return nil
}
