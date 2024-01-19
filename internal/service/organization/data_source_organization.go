package organization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasOrganizationRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
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
		},
	}
}

func dataSourceMongoDBAtlasOrganizationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)

	organization, _, err := conn.OrganizationsApi.GetOrganization(ctx, orgID).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organizations information: %s", err))
	}

	if err := d.Set("name", organization.GetName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `name`: %s", err))
	}

	if err := d.Set("is_deleted", organization.GetIsDeleted()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `is_deleted`: %s", err))
	}

	if err := d.Set("links", flattenOrganizationLinks(organization.GetLinks())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `is_deleted`: %s", err))
	}

	settings, _, err := conn.OrganizationsApi.GetOrganizationSettings(ctx, orgID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organization settings: %s", err))
	}
	if err := d.Set("api_access_list_required", settings.ApiAccessListRequired); err != nil {
		return diag.Errorf("error setting `api_access_list_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("multi_factor_auth_required", settings.MultiFactorAuthRequired); err != nil {
		return diag.Errorf("error setting `multi_factor_auth_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("restrict_employee_access", settings.RestrictEmployeeAccess); err != nil {
		return diag.Errorf("error setting `restrict_employee_access` for organization (%s): %s", orgID, err)
	}

	d.SetId(organization.GetId())

	return nil
}
