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
			"name": {
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
		},
	}
}

func dataSourceMongoDBAtlasProjectRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectName := d.Get("name").(string)

	project, _, err := conn.Projects.GetOneProjectByName(context.Background(), projectName)
	if err != nil {
		return fmt.Errorf("error getting project information: %s", err)
	}
	if err := d.Set("name", projectName); err != nil {
		return fmt.Errorf("error setting `name`: %s", err)
	}
	if err := d.Set("org_id", project.OrgID); err != nil {
		return fmt.Errorf("error setting `org_id` for project (%s): %s", d.Id(), err)
	}
	if err := d.Set("cluster_count", project.ClusterCount); err != nil {
		return fmt.Errorf("error setting `clusterCount` for project (%s): %s", d.Id(), err)
	}
	if err := d.Set("created", project.Created); err != nil {
		return fmt.Errorf("error setting `created` for project (%s): %s", d.Id(), err)
	}

	d.SetId(project.ID)
	return nil
}
