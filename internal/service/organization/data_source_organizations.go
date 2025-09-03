package organization

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: pluralDataSourceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"links": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rel": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"api_access_list_required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"multi_factor_auth_required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"restrict_employee_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"gen_ai_features_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"skip_default_alerts_settings": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"security_contact": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func pluralDataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	organizationOptions := &admin.ListOrganizationsApiParams{
		PageNum:      conversion.IntPtr(d.Get("page_num").(int)),
		ItemsPerPage: conversion.IntPtr(d.Get("items_per_page").(int)),
		Name:         conversion.Pointer(d.Get("name").(string)),
	}

	organizations, _, err := conn.OrganizationsApi.ListOrganizationsWithParams(ctx, organizationOptions).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organization information: %s", err))
	}
	results, err := flattenOrganizations(ctx, conn, organizations.GetResults())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("results", results); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	if err := d.Set("total_count", organizations.GetTotalCount()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `total_count`: %s", err))
	}

	d.SetId(id.UniqueId())
	return nil
}

func flattenOrganizations(ctx context.Context, conn *admin.APIClient, organizations []admin.AtlasOrganization) ([]map[string]any, error) {
	var results []map[string]any

	if len(organizations) == 0 {
		return results, nil
	}

	results = make([]map[string]any, len(organizations))

	for k, organization := range organizations {
		settings, _, err := conn.OrganizationsApi.GetOrganizationSettings(ctx, *organization.Id).Execute()
		if err != nil {
			return nil, fmt.Errorf("error getting organization settings (orgID: %s, org Name: %s): %s", organization.GetId(), organization.GetName(), err)
		}
		results[k] = map[string]any{
			"id":                           organization.Id,
			"name":                         organization.Name,
			"skip_default_alerts_settings": organization.SkipDefaultAlertsSettings,
			"is_deleted":                   organization.IsDeleted,
			"links":                        conversion.FlattenLinks(organization.GetLinks()),
			"api_access_list_required":     settings.ApiAccessListRequired,
			"multi_factor_auth_required":   settings.MultiFactorAuthRequired,
			"restrict_employee_access":     settings.RestrictEmployeeAccess,
			"gen_ai_features_enabled":      settings.GenAIFeaturesEnabled,
			"security_contact":             settings.SecurityContact,
		}
	}
	return results, nil
}
