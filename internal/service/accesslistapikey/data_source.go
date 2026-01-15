package accesslistapikey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_used": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_used_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)
	apiKeyID := d.Get("api_key_id").(string)
	ipAddress := d.Get("ip_address").(string)
	accessListAPIKey, _, err := connV2.ProgrammaticAPIKeysApi.GetOrgAccessEntry(ctx, orgID, ipAddress, apiKeyID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting access list api key information: %s", err))
	}

	if err := d.Set("cidr_block", accessListAPIKey.CidrBlock); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cidr_block`: %s", err))
	}

	if err := d.Set("last_used_address", accessListAPIKey.LastUsedAddress); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `last_used_address`: %s", err))
	}

	if err := d.Set("last_used", conversion.TimePtrToStringPtr(accessListAPIKey.LastUsed)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `last_used`: %s", err))
	}

	if err := d.Set("created", conversion.TimePtrToStringPtr(accessListAPIKey.Created)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `created`: %s", err))
	}

	if err := d.Set("access_count", accessListAPIKey.Count); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `access_count`: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}
