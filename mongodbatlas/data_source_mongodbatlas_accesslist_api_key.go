package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceMongoDBAtlasAccessListAPIKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasAccessListAPIKeyRead,
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

func dataSourceMongoDBAtlasAccessListAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	orgID := d.Get("org_id").(string)
	apiKeyID := d.Get("api_key_id").(string)
	ipAddress := d.Get("ip_address").(string)
	accessListAPIKey, _, err := conn.AccessListAPIKeys.Get(ctx, orgID, apiKeyID, ipAddress)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting access list api key information: %s", err))
	}

	if err := d.Set("cidr_block", accessListAPIKey.CidrBlock); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cidr_block`: %s", err))
	}

	if err := d.Set("last_used_address", accessListAPIKey.LastUsedAddress); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `last_used_address`: %s", err))
	}

	if err := d.Set("last_used", accessListAPIKey.LastUsed); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `last_used`: %s", err))
	}

	if err := d.Set("created", accessListAPIKey.Created); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `created`: %s", err))
	}

	if err := d.Set("access_count", accessListAPIKey.Count); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `access_count`: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}
