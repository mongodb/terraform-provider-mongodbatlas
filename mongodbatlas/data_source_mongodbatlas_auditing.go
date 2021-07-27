package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasAuditing() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasAuditingRead,
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

func dataSourceMongoDBAtlasAuditingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	auditing, _, err := conn.Auditing.Get(ctx, projectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	if err := d.Set("audit_authorization_success", auditing.AuditAuthorizationSuccess); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	if err := d.Set("audit_filter", auditing.AuditFilter); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	if err := d.Set("enabled", auditing.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	if err := d.Set("configuration_type", auditing.ConfigurationType); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	d.SetId(projectID)

	return nil
}
