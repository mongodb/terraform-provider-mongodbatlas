package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasProjectsRead,
		Schema: map[string]*schema.Schema{
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"org_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"teams": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"team_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role_names": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"api_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role_names": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"is_collect_database_specifics_statistics_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_data_explorer_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_extended_storage_sizes_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"is_performance_advisor_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_realtime_performance_panel_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_schema_advisor_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"region_usage_restrictions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"limits": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"current_usage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"default_limit": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"maximum_limit": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasProjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	client := meta.(*MongoDBClient)
	conn := client.Atlas
	var diags diag.Diagnostics

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	projects, _, err := conn.Projects.GetAllProjects(ctx, options)
	if err != nil {
		diagError := diag.Errorf("error getting projects information: %s", err)
		diags = append(diags, diagError...)
	}

	results, projectDiag := flattenProjects(ctx, client, projects.Results)
	if projectDiag != nil {
		diags = append(diags, projectDiag...)
	}

	if err := d.Set("results", results); err != nil {
		diagError := diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
		diags = append(diags, diagError...)
	}

	if err := d.Set("total_count", projects.TotalCount); err != nil {
		diagError := diag.FromErr(fmt.Errorf("error setting `name`: %s", err))
		diags = append(diags, diagError...)
	}

	d.SetId(id.UniqueId())

	return diags
}

func flattenProjects(ctx context.Context, client *MongoDBClient, projects []*matlas.Project) ([]map[string]interface{}, diag.Diagnostics) {
	conn := client.Atlas
	connV2 := client.AtlasV2

	var diags diag.Diagnostics

	var results []map[string]interface{}

	if len(projects) > 0 {
		results = make([]map[string]interface{}, len(projects))

		for k, project := range projects {
			teams, _, err := conn.Projects.GetProjectTeamsAssigned(ctx, project.ID)
			if err != nil {
				diagWarning := diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Error getting project's teams assigned",
					Detail:   fmt.Sprintf("Error getting project's teams assigned (%s): %s", project.ID, err),
				}
				diags = append(diags, diagWarning)
			}

			apiKeys, err := getProjectAPIKeys(ctx, conn, project.OrgID, project.ID)
			if err != nil {
				diagWarning := diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Error getting project's api keys",
					Detail:   fmt.Sprintf("Error getting project's api keys (%s): %s", project.ID, err),
				}
				diags = append(diags, diagWarning)
			}

			limits, _, err := connV2.ProjectsApi.ListProjectLimits(ctx, project.ID).Execute()
			if err != nil {
				diagWarning := diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Error getting project's limits",
					Detail:   fmt.Sprintf("Error getting project's limits (%s): %s", project.ID, err),
				}
				diags = append(diags, diagWarning)
			}

			projectSettings, _, err := conn.Projects.GetProjectSettings(ctx, project.ID)
			if err != nil {
				diagWarning := diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Error getting project's settings assigned",
					Detail:   fmt.Sprintf("Error getting project's settings assigned (%s): %s", project.ID, err),
				}
				diags = append(diags, diagWarning)
			}

			resultEntry := map[string]interface{}{
				"id":                        project.ID,
				"org_id":                    project.OrgID,
				"name":                      project.Name,
				"cluster_count":             project.ClusterCount,
				"created":                   project.Created,
				"region_usage_restrictions": project.RegionUsageRestrictions,
				"teams":                     flattenTeams(teams),
				"api_keys":                  flattenAPIKeys(apiKeys),
				"limits":                    flattenLimits(limits),
			}

			if projectSettings != nil {
				resultEntry["is_collect_database_specifics_statistics_enabled"] = projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled
				resultEntry["is_data_explorer_enabled"] = projectSettings.IsDataExplorerEnabled
				resultEntry["is_extended_storage_sizes_enabled"] = projectSettings.IsExtendedStorageSizesEnabled
				resultEntry["is_performance_advisor_enabled"] = projectSettings.IsPerformanceAdvisorEnabled
				resultEntry["is_realtime_performance_panel_enabled"] = projectSettings.IsRealtimePerformancePanelEnabled
				resultEntry["is_schema_advisor_enabled"] = projectSettings.IsSchemaAdvisorEnabled
			}

			results[k] = resultEntry
		}
	}

	return results, diags
}
