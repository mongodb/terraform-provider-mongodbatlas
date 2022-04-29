package mongodbatlas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasPrivateEndpointRegionalMode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasPrivateEndpointRegionalModeRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasPrivateEndpointRegionalModeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	setting, _, err := conn.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(ctx, projectID)
	if err != nil {
		// case 404
		// deleted in the backend case
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
