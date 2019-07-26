package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func resourceMongoDBAtlasProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasProjectCreate,
		Read:   resourceMongoDBAtlasProjectRead,
		Delete: resourceMongoDBAtlasProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceMongoDBAtlasProjectCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	projectReq := &matlas.Project{
		OrgID: d.Get("org_id").(string),
		Name:  d.Get("name").(string),
	}

	projectRes, _, err := conn.Projects.Create(context.Background(), projectReq)
	if err != nil {
		return fmt.Errorf("error creating project: %s", err)
	}

	d.SetId(projectRes.ID)
	return resourceMongoDBAtlasProjectRead(d, meta)
}

func resourceMongoDBAtlasProjectRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Id()

	projectRes, _, err := conn.Projects.GetOneProject(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("error getting project information: %s", err)
	}

	if err := d.Set("name", projectRes.Name); err != nil {
		return fmt.Errorf("error setting `name` for project (%s): %s", d.Id(), err)
	}
	if err := d.Set("org_id", projectRes.OrgID); err != nil {
		return fmt.Errorf("error setting `org_id` for project (%s): %s", d.Id(), err)
	}
	if err := d.Set("cluster_count", projectRes.ClusterCount); err != nil {
		return fmt.Errorf("error setting `clusterCount` for project (%s): %s", d.Id(), err)
	}
	if err := d.Set("created", projectRes.Created); err != nil {
		return fmt.Errorf("error setting `created` for project (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceMongoDBAtlasProjectDelete(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Id()

	_, err := conn.Projects.Delete(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("error deleting project (%s): %s", projectID, err)
	}
	return nil
}
