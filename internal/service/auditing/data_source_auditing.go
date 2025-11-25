package auditing

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	auditing, _, err := connV2.AuditingApi.GetGroupAuditLog(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	if err := d.Set("audit_authorization_success", auditing.GetAuditAuthorizationSuccess()); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	if err := d.Set("audit_filter", auditing.GetAuditFilter()); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	if err := d.Set("enabled", auditing.GetEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	if err := d.Set("configuration_type", auditing.GetConfigurationType()); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, projectID, err))
	}

	d.SetId(projectID)

	return nil
}
