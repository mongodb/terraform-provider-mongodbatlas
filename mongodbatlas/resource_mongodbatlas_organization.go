package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorOrganizationCreate  = "error creating MongoDB Organization: %s"
	errorOrganizationRead    = "error reading MongoDB Organization (%s): %s"
	errorOrganizationUpate   = "error updating MongoDB Organization (%s): %s"
	errorOrganizationDelete  = "error deleting MongoDB Organization (%s): %s"
	errorOrganizationSetting = "error setting `%s` for organization (%s): %s"
)

func resourceMongoDBAtlasOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasOrganizationCreate,
		Read:   resourceMongoDBAtlasOrganizationRead,
		Update: resourceMongoDBAtlasOrganizationUpdate,
		Delete: resourceMongoDBAtlasOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceMongoDBAtlasOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	organization, _, err := conn.Organizations.Create(context.Background(), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf(errorOrganizationCreate, err)
	}

	d.SetId(organization.ID)
	return resourceMongoDBAtlasOrganizationRead(d, meta)
}

func resourceMongoDBAtlasOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	organizationID := d.Id()

	_, _, err := conn.Organizations.GetOneOrganization(context.Background(), organizationID)
	if err != nil {
		return fmt.Errorf(errorOrganizationRead, organizationID, err)
	}

	return nil
}

func resourceMongoDBAtlasOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	if d.HasChange("name") {
		req := &matlas.Organization{
			ID:   d.Id(),
			Name: d.Get("name").(string),
		}

		_, _, err := conn.Organizations.UpdateOrganizationName(context.Background(), req)
		if err != nil {
			return fmt.Errorf(errorOrganizationUpate, d.Id(), err)
		}
	}

	d.SetId(d.Id())
	return resourceMongoDBAtlasOrganizationRead(d, meta)
}

func resourceMongoDBAtlasOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	_, err := conn.Organizations.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf(errorOrganizationDelete, d.Id(), err)
	}

	return nil
}
