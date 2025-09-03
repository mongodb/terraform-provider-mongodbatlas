package projectapikey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: pluralDataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
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
						"private_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_assignment": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"role_names": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func pluralDataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	params := &admin.ListProjectApiKeysApiParams{
		GroupId:      d.Get("project_id").(string),
		PageNum:      conversion.IntPtr(d.Get("page_num").(int)),
		ItemsPerPage: conversion.IntPtr(d.Get("items_per_page").(int)),
	}
	apiKeys, _, err := connV2.ProgrammaticAPIKeysApi.ListProjectApiKeysWithParams(ctx, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api keys information: %s", err))
	}

	results, err := flattenProjectAPIKeys(ctx, connV2, apiKeys.GetResults())
	if err != nil {
		diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	if err := d.Set("results", results); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenProjectAPIKeys(ctx context.Context, connV2 *admin.APIClient, apiKeys []admin.ApiKeyUserDetails) ([]map[string]any, error) {
	var results []map[string]any

	if len(apiKeys) == 0 {
		return nil, nil
	}

	results = make([]map[string]any, len(apiKeys))
	for k, apiKey := range apiKeys {
		results[k] = map[string]any{
			"api_key_id":  apiKey.GetId(),
			"description": apiKey.GetDesc(),
			"public_key":  apiKey.GetPublicKey(),
			"private_key": apiKey.GetPrivateKey(),
		}

		details, _, err := getKeyDetails(ctx, connV2, apiKey.GetId())
		if err != nil {
			return nil, err
		}
		results[k]["project_assignment"] = flattenProjectAssignments(details.GetRoles())
	}
	return results, nil
}
