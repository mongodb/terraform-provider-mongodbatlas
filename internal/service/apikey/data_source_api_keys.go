package apikey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_names": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
	params := &admin.ListApiKeysApiParams{
		PageNum:      conversion.IntPtr(d.Get("page_num").(int)),
		ItemsPerPage: conversion.IntPtr(d.Get("items_per_page").(int)),
		OrgId:        orgID,
	}

	apiKeys, _, err := connV2.ProgrammaticAPIKeysApi.ListApiKeysWithParams(ctx, params).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api keys information: %s", err))
	}

	if err := d.Set("results", flattenOrgAPIKeys(ctx, orgID, apiKeys.GetResults())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	d.SetId(id.UniqueId())
	return nil
}

func flattenOrgAPIKeys(ctx context.Context, orgID string, apiKeys []admin.ApiKeyUserDetails) []map[string]any {
	results := make([]map[string]any, len(apiKeys))
	for k, apiKey := range apiKeys {
		results[k] = map[string]any{
			"api_key_id":  apiKey.GetId(),
			"description": apiKey.GetDesc(),
			"public_key":  apiKey.GetPublicKey(),
			"role_names":  flattenOrgAPIKeyRoles(orgID, apiKey.GetRoles()),
		}
	}
	return results
}
