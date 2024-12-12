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
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	dnsResp, _, err := connV2.AWSClustersDNSApi.GetAwsCustomDns(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorRead, err))
	}
	if err := d.Set("enabled", dnsResp.GetEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSetting, "enabled", projectID, err))
	}
	d.SetId(projectID)
	return nil
}
