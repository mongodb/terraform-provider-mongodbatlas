package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasAuditing() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasAuditingRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"audit_authorization_success": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"audit_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"configuration_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasAuditingRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)

	auditing, _, err := conn.Auditings.Get(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf(errorAuditingRead, projectID, err)
	}

	if err := d.Set("audit_authorization_success", auditing.AuditAuthorizationSuccess); err != nil {
		return fmt.Errorf(errorAuditingRead, projectID, err)
	}
	if err := d.Set("audit_filter", auditing.AuditFilter); err != nil {
		return fmt.Errorf(errorAuditingRead, projectID, err)
	}
	if err := d.Set("enabled", auditing.Enabled); err != nil {
		return fmt.Errorf(errorAuditingRead, projectID, err)
	}
	if err := d.Set("configuration_type", auditing.ConfigurationType); err != nil {
		return fmt.Errorf(errorAuditingRead, projectID, err)
	}

	d.SetId(projectID)
	return nil
}
