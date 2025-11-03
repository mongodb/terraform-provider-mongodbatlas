package privateendpointregionalmode

import (
	"context"
	"strings"

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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	setting, _, err := conn.PrivateEndpointServicesApi.GetRegionalEndpointMode(ctx, projectID).Execute()
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting private endpoint regional mode: %s", err)
	}

	if err := d.Set("enabled", setting.Enabled); err != nil {
		return diag.Errorf("error setting `enabled` for enabled (%s): %s", d.Id(), err)
	}

	d.SetId(projectID)

	return nil
}
