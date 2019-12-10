package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasAuditingRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasAuditingRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	orgID := d.Get("org_id").(string)

	org, _, err := conn.Organizations.GetOneOrganization(context.Background(), orgID)
	if err != nil {
		return fmt.Errorf(errorOrganizationRead, orgID, err)
	}

	if err := d.Set("name", org.Name); err != nil {
		return fmt.Errorf(errorOrganizationSetting, "name", orgID, err)
	}

	d.SetId(orgID)
	return nil
}
