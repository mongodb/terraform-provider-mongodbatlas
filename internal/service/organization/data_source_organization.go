package organization

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
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
			"users": dsschema.DSOrgUsersSchema(),
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
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)

	organization, _, err := conn.OrganizationsApi.GetOrganization(ctx, orgID).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organizations information: %s", err))
	}

	if err := d.Set("name", organization.GetName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `name`: %s", err))
	}

	if err := d.Set("skip_default_alerts_settings", organization.GetSkipDefaultAlertsSettings()); err != nil {
		return diag.Errorf("error setting `skip_default_alerts_settings`: %s", err)
	}

	if err := d.Set("is_deleted", organization.GetIsDeleted()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `is_deleted`: %s", err))
	}

	if err := d.Set("links", conversion.FlattenLinks(organization.GetLinks())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `is_deleted`: %s", err))
	}

	users, err := listAllOrganizationUsers(ctx, orgID, conn)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organization users: %s", err))
	}
	if err := d.Set("users", conversion.FlattenUsers(users)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `users`: %s", err))
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
	if err := d.Set("gen_ai_features_enabled", settings.GenAIFeaturesEnabled); err != nil {
		return diag.Errorf("error setting `gen_ai_features_enabled` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("security_contact", settings.SecurityContact); err != nil {
		return diag.Errorf("error setting `security_contact` for organization (%s): %s", orgID, err)
	}

	d.SetId(organization.GetId())

	return nil
}

func listAllOrganizationUsers(ctx context.Context, orgID string, conn *admin.APIClient) ([]admin.OrgUserResponse, error) {
	return dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.OrgUserResponse], *http.Response, error) {
		request := conn.MongoDBCloudUsersApi.ListOrganizationUsers(ctx, orgID)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
}
