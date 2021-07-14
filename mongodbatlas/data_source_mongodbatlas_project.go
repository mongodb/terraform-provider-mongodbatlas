package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasProjectRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"project_id"},
			},
			"org_id": {
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
		},
	}
}

func dataSourceMongoDBAtlasProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	projectID, projectIDOk := d.GetOk("project_id")
	name, nameOk := d.GetOk("name")

	if !projectIDOk && !nameOk {
		return diag.FromErr(errors.New("either project_id or name must be configured"))
	}

	var (
		err     error
		project *matlas.Project
	)

	if projectIDOk {
		project, _, err = conn.Projects.GetOneProject(ctx, projectID.(string))
	} else {
		project, _, err = conn.Projects.GetOneProjectByName(ctx, name.(string))
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorProjectRead, projectID, err))
	}

	teams, _, err := conn.Projects.GetProjectTeamsAssigned(ctx, project.ID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting project's teams assigned (%s): %s", projectID, err))
	}

	if err := d.Set("org_id", project.OrgID); err != nil {
		return diag.FromErr(fmt.Errorf(errorProjectSetting, `org_id`, project.ID, err))
	}

	if err := d.Set("cluster_count", project.ClusterCount); err != nil {
		return diag.FromErr(fmt.Errorf(errorProjectSetting, `clusterCount`, project.ID, err))
	}

	if err := d.Set("created", project.Created); err != nil {
		return diag.FromErr(fmt.Errorf(errorProjectSetting, `created`, project.ID, err))
	}

	if err := d.Set("teams", flattenTeams(teams)); err != nil {
		return diag.FromErr(fmt.Errorf(errorProjectSetting, `teams`, project.ID, err))
	}

	d.SetId(project.ID)

	return nil
}
