package accesslistapikey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)
	apiKeyID := d.Get("api_key_id").(string)
	params := &admin.ListOrgAccessEntriesApiParams{
		PageNum:      conversion.IntPtr(d.Get("page_num").(int)),
		ItemsPerPage: conversion.IntPtr(d.Get("items_per_page").(int)),
		OrgId:        orgID,
		ApiUserId:    apiKeyID,
	}
	accessListAPIKeys, _, err := connV2.ProgrammaticAPIKeysApi.ListOrgAccessEntriesWithParams(ctx, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting access list api keys information: %s", err))
	}

	if err := d.Set("results", flattenAccessListAPIKeys(ctx, orgID, accessListAPIKeys.GetResults())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenAccessListAPIKeys(ctx context.Context, orgID string, list []admin.UserAccessListResponse) []map[string]any {
	results := make([]map[string]any, len(list))
	for k, elm := range list {
		results[k] = map[string]any{
			"ip_address":        elm.IpAddress,
			"cidr_block":        elm.CidrBlock,
			"created":           conversion.TimePtrToStringPtr(elm.Created),
			"access_count":      elm.Count,
			"last_used":         conversion.TimePtrToStringPtr(elm.LastUsed),
			"last_used_address": elm.LastUsedAddress,
		}
	}
	return results
}
