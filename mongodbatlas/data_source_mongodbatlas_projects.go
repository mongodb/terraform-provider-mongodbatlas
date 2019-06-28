package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasProjectsRead,
		Schema: map[string]*schema.Schema{
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

func dataSourceMongoDBAtlasProjectsRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	projects, _, err := conn.Projects.GetAllProjects(context.Background())
	if err != nil {
		return fmt.Errorf("error getting projects information: %s", err)
	}
	if err := d.Set("results", flattenProjects(projects.Results)); err != nil {
		return fmt.Errorf("error setting `results`: %s", err)
	}
	if err := d.Set("total_count", projects.TotalCount); err != nil {
		return fmt.Errorf("error setting `name`: %s", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenProjects(projects []*matlas.Project) []map[string]interface{} {
	var results []map[string]interface{}

	if len(projects) > 0 {
		results = make([]map[string]interface{}, len(projects))

		for k, project := range projects {
			results[k] = map[string]interface{}{
				"id":            project.ID,
				"org_id":        project.OrgID,
				"name":          project.Name,
				"cluster_count": project.ClusterCount,
				"created":       project.Created,
			}
		}
	}
	return results
}
