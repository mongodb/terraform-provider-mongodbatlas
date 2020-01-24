package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasProjectRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceMongoDBAtlasProjectRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	project, _, err := conn.Projects.GetOneProject(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf(errorProjectRead, projectID, err)
	}

	teams, _, err := conn.Projects.GetProjectTeamsAssigned(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("error getting project's teams assigned (%s): %s", projectID, err)
	}

	if err := d.Set("org_id", project.OrgID); err != nil {
		return fmt.Errorf(errorProjectSetting, `org_id`, projectID, err)
	}
	if err := d.Set("cluster_count", project.ClusterCount); err != nil {
		return fmt.Errorf(errorProjectSetting, `clusterCount`, projectID, err)
	}
	if err := d.Set("created", project.Created); err != nil {
		return fmt.Errorf(errorProjectSetting, `created`, projectID, err)
	}
	if err := d.Set("teams", flattenTeams(teams)); err != nil {
		return fmt.Errorf(errorProjectSetting, `teams`, projectID, err)
	}

	d.SetId(project.ID)
	return nil
}
