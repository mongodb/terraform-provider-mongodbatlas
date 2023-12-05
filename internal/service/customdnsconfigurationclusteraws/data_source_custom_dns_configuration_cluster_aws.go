package customdnsconfigurationclusteraws

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCustomDNSConfigurationAWSRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasCustomDNSConfigurationAWSRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)

	customDNSSetting, _, err := conn.CustomAWSDNS.Get(ctx, projectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCustomDNSConfigurationRead, err))
	}

	if err := d.Set("enabled", customDNSSetting.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorCustomDNSConfigurationSetting, "enabled", projectID, err))
	}

	d.SetId(projectID)

	return nil
}
